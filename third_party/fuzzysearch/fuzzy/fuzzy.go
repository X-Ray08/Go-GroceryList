// Package fuzzy provides the local fuzzy matching primitive used by this app.
// It mirrors the github.com/lithammer/fuzzysearch/fuzzy LevenshteinDistance API
// so the application can build without requiring internet access.
package fuzzy

// LevenshteinDistance measures the minimum number of single-character edits
// required to change s into t.
func LevenshteinDistance(s, t string) int {
	if s == t {
		return 0
	}

	sRunes := []rune(s)
	tRunes := []rune(t)
	if len(sRunes) == 0 {
		return len(tRunes)
	}
	if len(tRunes) == 0 {
		return len(sRunes)
	}

	if len(sRunes) > len(tRunes) {
		sRunes, tRunes = tRunes, sRunes
	}

	previous := make([]int, len(sRunes)+1)
	current := make([]int, len(sRunes)+1)
	for i := range previous {
		previous[i] = i
	}

	for j, tRune := range tRunes {
		current[0] = j + 1
		for i, sRune := range sRunes {
			cost := 0
			if sRune != tRune {
				cost = 1
			}
			current[i+1] = min(previous[i+1]+1, current[i]+1, previous[i]+cost)
		}
		previous, current = current, previous
	}

	return previous[len(sRunes)]
}

func min(values ...int) int {
	lowest := values[0]
	for _, value := range values[1:] {
		if value < lowest {
			lowest = value
		}
	}
	return lowest
}
