package model

// Item is a single grocery list entry.
type Item struct {
	ID       int
	Name     string
	Category string
	Bought   bool
}
