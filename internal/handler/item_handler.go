package handler

import (
	"net/http"
	"strconv"

	"go-grocerylist/internal/model"
	"go-grocerylist/internal/repository"
	"go-grocerylist/internal/service"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct{ service *service.ItemService }

type pageData struct {
	Groups        map[string][]model.Item
	Ordered       []string
	Total         int
	Bought        int
	Progress      string
	Error         string
	Filter        string
	HasItems      bool
	HasFilterHits bool
}

func NewItemHandler(service *service.ItemService) *ItemHandler { return &ItemHandler{service: service} }

func (h *ItemHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.Index)
	r.GET("/items", h.ListItems)
	r.GET("/items/filter/:type", h.FilterItems)
	r.POST("/items", h.CreateItem)
	r.PUT("/items/:id/toggle", h.ToggleItem)
}

func (h *ItemHandler) Index(c *gin.Context)     { h.renderList(c, http.StatusOK, "all", "") }
func (h *ItemHandler) ListItems(c *gin.Context) { h.renderList(c, http.StatusOK, "all", "") }
func (h *ItemHandler) FilterItems(c *gin.Context) {
	h.renderList(c, http.StatusOK, c.Param("type"), "")
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	if _, err := h.service.AddItem(c.PostForm("name")); err != nil {
		h.renderList(c, http.StatusBadRequest, "all", err.Error())
		return
	}
	h.renderList(c, http.StatusOK, "all", "")
}

func (h *ItemHandler) ToggleItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.renderList(c, http.StatusBadRequest, "all", "invalid item id")
		return
	}
	if _, err = h.service.ToggleItem(id); err != nil {
		if err == repository.ErrNotFound {
			h.renderList(c, http.StatusNotFound, "all", "item not found")
			return
		}
		h.renderList(c, http.StatusInternalServerError, "all", "failed to toggle item")
		return
	}
	h.renderList(c, http.StatusOK, "all", "")
}

func (h *ItemHandler) renderList(c *gin.Context, status int, filter, errMsg string) {
	data, err := h.buildPageData(filter, errMsg)
	if err != nil {
		c.HTML(http.StatusBadRequest, "item_list.html", pageData{Error: err.Error(), Filter: "all"})
		return
	}
	c.HTML(status, "item_list.html", data)
}

func (h *ItemHandler) buildPageData(filter, errMsg string) (pageData, error) {
	allItems, err := h.service.ListItems("all")
	if err != nil {
		return pageData{}, err
	}
	items, err := h.service.ListItems(filter)
	if err != nil {
		return pageData{}, err
	}
	groups := map[string][]model.Item{}
	ordered, seen := []string{}, map[string]bool{}
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
	return pageData{Groups: groups, Ordered: ordered, Total: len(items), Bought: bought, Progress: strconv.Itoa(bought) + " of " + strconv.Itoa(len(items)) + " items bought", Error: errMsg, Filter: filter, HasItems: len(allItems) > 0, HasFilterHits: len(items) > 0}, nil
}
