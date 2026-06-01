package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go-grocerylist/internal/model"
)

var ErrNotFound = errors.New("item not found")

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Add(item model.Item) (model.Item, error) {
	result, err := r.db.Exec(`INSERT INTO items(name, category, bought) VALUES (?, ?, ?)`, item.Name, item.Category, item.Bought)
	if err != nil {
		log.Printf("database error creating item: %v", err)
		return model.Item{}, fmt.Errorf("insert item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Item{}, fmt.Errorf("get inserted id: %w", err)
	}

	return r.GetByID(int(id))
}

func (r *ItemRepository) Toggle(id int) (model.Item, error) {
	res, err := r.db.Exec(`UPDATE items SET bought = CASE WHEN bought = 1 THEN 0 ELSE 1 END WHERE id = ?`, id)
	if err != nil {
		log.Printf("database error toggling item: %v", err)
		return model.Item{}, fmt.Errorf("toggle item: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return model.Item{}, fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return model.Item{}, ErrNotFound
	}

	return r.GetByID(id)
}

func (r *ItemRepository) GetByID(id int) (model.Item, error) {
	row := r.db.QueryRow(`SELECT id, name, category, bought, created_at FROM items WHERE id = ?`, id)
	var item model.Item
	if err := row.Scan(&item.ID, &item.Name, &item.Category, &item.Bought, &item.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Item{}, ErrNotFound
		}
		return model.Item{}, fmt.Errorf("query item by id: %w", err)
	}
	return item, nil
}

func (r *ItemRepository) List(filter string) ([]model.Item, error) {
	query := `SELECT id, name, category, bought, created_at FROM items`
	switch filter {
	case "bought":
		query += ` WHERE bought = 1`
	case "unbought":
		query += ` WHERE bought = 0`
	}
	query += ` ORDER BY bought ASC, LOWER(name) ASC, created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("database error listing items: %v", err)
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	items := make([]model.Item, 0)
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Bought, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate items: %w", err)
	}

	return items, nil
}
