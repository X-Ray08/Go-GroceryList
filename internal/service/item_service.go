package service

import (
	"errors"
	"log"
	"strings"

	"go-grocerylist/internal/model"
	"go-grocerylist/internal/repository"
)

var (
	ErrInvalidName   = errors.New("item name cannot be empty")
	ErrNameTooLong   = errors.New("item name must be 80 characters or less")
	ErrInvalidFilter = errors.New("invalid filter")
)

type ItemService struct {
	repo *repository.ItemRepository
}

func NewItemService(repo *repository.ItemRepository) *ItemService { return &ItemService{repo: repo} }

func (s *ItemService) AddItem(name string) (model.Item, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return model.Item{}, ErrInvalidName
	}
	if len([]rune(cleanName)) > 80 {
		return model.Item{}, ErrNameTooLong
	}

	created, err := s.repo.Add(model.Item{Name: cleanName, Category: categorize(cleanName), Bought: false})
	if err != nil {
		return model.Item{}, err
	}
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

func (s *ItemService) ListItems(filter string) ([]model.Item, error) {
	if filter == "" {
		filter = "all"
	}
	if filter != "all" && filter != "bought" && filter != "unbought" {
		return nil, ErrInvalidFilter
	}
	return s.repo.List(filter)
}

func categorize(name string) string {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "milk") || strings.Contains(lower, "cheese") || strings.Contains(lower, "yogurt") {
		return "Dairy"
	}
	if strings.Contains(lower, "apple") || strings.Contains(lower, "banana") {
		return "Fruit"
	}
	if strings.Contains(lower, "chicken") || strings.Contains(lower, "beef") {
		return "Meat"
	}
	if strings.Contains(lower, "bread") {
		return "Bakery"
	}
	return "Other"
}
