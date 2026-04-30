package service

import (
	"errors"
	"log"
	"strings"

	"go-grocerylist/internal/model"
	"go-grocerylist/internal/repository"
)

var ErrInvalidName = errors.New("item name cannot be empty")

type ItemService struct {
	repo *repository.ItemRepository
}

func NewItemService(repo *repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) AddItem(name string) (model.Item, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return model.Item{}, ErrInvalidName
	}

	item := model.Item{
		Name:     cleanName,
		Category: categorize(cleanName),
		Bought:   false,
	}
	created := s.repo.Add(item)
	log.Printf("created item id=%d name=%q category=%q", created.ID, created.Name, created.Category)
	return created, nil
}

func (s *ItemService) ToggleItem(id int) (model.Item, error) {
	item, err := s.repo.Toggle(id)
	if err != nil {
		return model.Item{}, err
	}
	log.Printf("toggled item id=%d bought=%t", item.ID, item.Bought)
	return item, nil
}

func (s *ItemService) ListItems() []model.Item {
	return s.repo.ListSorted()
}

func categorize(name string) string {
	lower := strings.ToLower(name)

	rules := map[string][]string{
		"Dairy":      {"milk", "cheese", "yogurt", "butter", "cream"},
		"Fruit":      {"apple", "banana", "orange", "berry", "grape", "mango"},
		"Vegetables": {"carrot", "broccoli", "spinach", "lettuce", "tomato", "onion"},
		"Bakery":     {"bread", "bun", "bagel", "croissant"},
		"Meat":       {"chicken", "beef", "pork", "fish", "turkey"},
		"Beverages":  {"juice", "soda", "coffee", "tea", "water"},
	}

	for category, keywords := range rules {
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				return category
			}
		}
	}
	return "Other"
}
