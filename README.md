# Smart Grocery List

A grocery list app I built with Go, Gin, HTMX and SQLite. When you add an item, it gets sorted into a store category automatically (Dairy, Produce, Bakery...). Typos are handled too, so "bananna" still lands in Fruit.

## Run it

You need Go 1.22+ and `gcc` (the SQLite driver uses cgo). From the repo root:

```bash
go run ./cmd
```

Then open http://localhost:8080. Set `PORT` to use a different port.

Tests: `go test ./...`

## How it works

- `cmd/main.go` sets everything up
- `internal/` holds the handlers, service, repository and SQLite code
- `web/` has the HTML templates, CSS and HTMX
- `data/products.json` maps product names to categories

Categorizing tries an exact match first, then a partial match, then a fuzzy match. If nothing fits, the item goes under "Other". To add a product, put it in `products.json` and restart.

Items are saved in `grocery.db`. Delete that file to start fresh.
