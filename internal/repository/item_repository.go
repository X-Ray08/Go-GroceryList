package repository

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"go-grocerylist/internal/model"
)

var ErrNotFound = errors.New("item not found")

// ItemRepository stores items in-memory using a mutex for thread safety.
type ItemRepository struct {
	mu     sync.RWMutex
	items  []model.Item
	nextID int
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{nextID: 1}
}

func (r *ItemRepository) Add(item model.Item) model.Item {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = r.nextID
	r.nextID++
	r.items = append(r.items, item)
	return item
}

func (r *ItemRepository) Toggle(id int) (model.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.items {
		if r.items[i].ID == id {
			r.items[i].Bought = !r.items[i].Bought
			return r.items[i], nil
		}
	}

	return model.Item{}, ErrNotFound
}

func (r *ItemRepository) ListSorted() []model.Item {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]model.Item, len(r.items))
	copy(copied, r.items)

	// Unbought items first, then by name.
	sort.SliceStable(copied, func(i, j int) bool {
		if copied[i].Bought != copied[j].Bought {
			return !copied[i].Bought
		}
		return strings.ToLower(copied[i].Name) < strings.ToLower(copied[j].Name)
	})

	return copied
}
