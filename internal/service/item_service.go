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
	repo        *repository.ItemRepository
	categorizer *Categorizer
}

func NewItemService(repo *repository.ItemRepository, categorizer *Categorizer) *ItemService {
	return &ItemService{repo: repo, categorizer: categorizer}
}

func (s *ItemService) AddItem(name string) (model.Item, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return model.Item{}, ErrInvalidName
	}
	if len([]rune(cleanName)) > 80 {
		return model.Item{}, ErrNameTooLong
	}

	created, err := s.repo.Add(model.Item{Name: cleanName, Category: s.categorizer.Categorize(cleanName), Bought: false})
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

func (s *ItemService) DeleteItem(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	log.Printf("deleted item id=%d", id)
	return nil
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
