package main

import (
	"log"
	"os"

	"html/template"

	"github.com/gin-gonic/gin"
	"go-grocerylist/internal/handler"
	"go-grocerylist/internal/repository"
	"go-grocerylist/internal/service"
)

func main() {
	r := gin.Default()

	r.Static("/static", "web/static")
	r.SetFuncMap(template.FuncMap{
		"dict": func(values ...any) map[string]any {
			d := make(map[string]any)
			for i := 0; i < len(values)-1; i += 2 {
				key, _ := values[i].(string)
				d[key] = values[i+1]
			}
			return d
		},
	})
	r.LoadHTMLGlob("web/templates/*.html")

	repo := repository.NewItemRepository()
	svc := service.NewItemService(repo)
	h := handler.NewItemHandler(svc)
	h.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("starting Smart Grocery List on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
