package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

const otherCategory = "Other"

// Categorizer assigns grocery categories from a local product dictionary.
// It deliberately avoids network calls and AI APIs so categorization remains
// fast, deterministic, and available offline.
type Categorizer struct {
	products map[string]string
	names    []string
}

// NewCategorizer loads product names and categories from a JSON dictionary.
func NewCategorizer(path string) (*Categorizer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open product dictionary: %w", err)
	}
	defer file.Close()

	rawProducts := map[string]string{}
	if err := json.NewDecoder(file).Decode(&rawProducts); err != nil {
		return nil, fmt.Errorf("decode product dictionary: %w", err)
	}
	if len(rawProducts) == 0 {
		return nil, fmt.Errorf("product dictionary is empty")
	}

	products := make(map[string]string, len(rawProducts))
	for name, category := range rawProducts {
		cleanName := normalize(name)
		cleanCategory := strings.TrimSpace(category)
		if cleanName == "" || cleanCategory == "" {
			continue
		}
		products[cleanName] = cleanCategory
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("product dictionary has no usable entries")
	}

	names := make([]string, 0, len(products))
	for name := range products {
		names = append(names, name)
	}
	// Check longer names first so "vanilla ice cream" is categorized as Frozen
	// before the shorter "cream" entry can match it as Dairy.
	sort.Slice(names, func(i, j int) bool {
		if len(names[i]) == len(names[j]) {
			return names[i] < names[j]
		}
		return len(names[i]) > len(names[j])
	})

	log.Printf("loaded %d grocery products from %s", len(products), path)
	return &Categorizer{products: products, names: names}, nil
}

// Categorize uses exact, partial, then typo-tolerant fuzzy matching before
// falling back to Other.
func (c *Categorizer) Categorize(itemName string) string {
	query := normalize(itemName)
	if query == "" {
		log.Printf("Categorized %q as %s using fallback to Other.", itemName, otherCategory)
		return otherCategory
	}

	if category, ok := c.products[query]; ok {
		log.Printf("Categorized %q as %s using exact match.", itemName, category)
		return category
	}

	for _, productName := range c.names {
		if strings.Contains(" "+query+" ", " "+productName+" ") {
			category := c.products[productName]
			log.Printf("Categorized %q as %s using partial match.", itemName, category)
			return category
		}
	}

	if productName, category, ok := c.bestFuzzyMatch(query); ok {
		log.Printf("Categorized %q as %s using fuzzy match for %q.", itemName, category, productName)
		return category
	}

	log.Printf("Categorized %q as %s using fallback to Other.", itemName, otherCategory)
	return otherCategory
}

func (c *Categorizer) bestFuzzyMatch(query string) (string, string, bool) {
	bestName := ""
	bestDistance := 0
	bestLength := 1

	for _, candidate := range fuzzyInputs(query) {
		for _, productName := range c.names {
			if !sameFirstRune(candidate, productName) {
				continue
			}
			distance := fuzzy.LevenshteinDistance(candidate, productName)
			if bestName == "" || distance < bestDistance || (distance == bestDistance && len(productName) > bestLength) {
				bestName = productName
				bestDistance = distance
				bestLength = len(productName)
			}
		}
	}

	if bestName == "" || !reasonableFuzzyMatch(bestDistance, bestLength) {
		return "", "", false
	}
	return bestName, c.products[bestName], true
}

func fuzzyInputs(query string) []string {
	inputs := []string{query}
	words := strings.Fields(query)
	for i := range words {
		inputs = append(inputs, words[i])
		for j := i + 1; j < len(words); j++ {
			inputs = append(inputs, strings.Join(words[i:j+1], " "))
		}
	}
	return inputs
}

func sameFirstRune(left, right string) bool {
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	return len(leftRunes) > 0 && len(rightRunes) > 0 && leftRunes[0] == rightRunes[0]
}

func reasonableFuzzyMatch(distance, productLength int) bool {
	allowedDistance := 1
	if productLength >= 5 {
		allowedDistance = 2
	}
	if productLength >= 9 {
		allowedDistance = 3
	}
	return distance <= allowedDistance && float64(distance)/float64(productLength) <= 0.35
}

func normalize(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(value))), " ")
}
