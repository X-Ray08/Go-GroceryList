package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-grocerylist/internal/model"
	"go-grocerylist/internal/repository"
	"go-grocerylist/internal/service"
)

type ItemHandler struct {
	service *service.ItemService
}

type pageData struct {
	Groups   map[string][]model.Item
	Ordered  []string
	Total    int
	Bought   int
	Progress string
	Error    string
}

func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.Index)
	r.GET("/items", h.ListItems)
	r.POST("/items", h.CreateItem)
	r.PUT("/items/:id/toggle", h.ToggleItem)
}

func (h *ItemHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", h.buildPageData(""))
}

func (h *ItemHandler) ListItems(c *gin.Context) {
	c.HTML(http.StatusOK, "item_list.html", h.buildPageData(""))
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	name := c.PostForm("name")
	_, err := h.service.AddItem(name)
	if err != nil {
		c.HTML(http.StatusBadRequest, "item_list.html", h.buildPageData(err.Error()))
		return
	}
	c.HTML(http.StatusOK, "item_list.html", h.buildPageData(""))
}

func (h *ItemHandler) ToggleItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "item_list.html", h.buildPageData("invalid item id"))
		return
	}

	_, err = h.service.ToggleItem(id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.HTML(http.StatusNotFound, "item_list.html", h.buildPageData("item not found"))
			return
		}
		c.HTML(http.StatusInternalServerError, "item_list.html", h.buildPageData("failed to toggle item"))
		return
	}

	c.HTML(http.StatusOK, "item_list.html", h.buildPageData(""))
}

func (h *ItemHandler) buildPageData(errMsg string) pageData {
	items := h.service.ListItems()
	groups := map[string][]model.Item{}
	ordered := []string{}
	seen := map[string]bool{}
	bought := 0

	for _, item := range items {
		groups[item.Category] = append(groups[item.Category], item)
		if !seen[item.Category] {
			ordered = append(ordered, item.Category)
			seen[item.Category] = true
		}
		if item.Bought {
			bought++
		}
	}

	return pageData{
		Groups:   groups,
		Ordered:  ordered,
		Total:    len(items),
		Bought:   bought,
		Progress: strconv.Itoa(bought) + " of " + strconv.Itoa(len(items)) + " items bought",
		Error:    errMsg,
	}
}
