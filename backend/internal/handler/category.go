package handler

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	*Handler
	svc service.CategoryServiceInterface
}

func NewCategoryHandler(h *Handler, svc service.CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{Handler: h, svc: svc}
}
func categoryJSON(item model.Category) v1.Category {
	out := v1.Category{ID: item.ID, ParentID: item.ParentID, Name: item.Name, Sort: item.Sort, Enabled: item.Enabled}
	for _, child := range item.Children {
		out.Children = append(out.Children, categoryJSON(child))
	}
	return out
}

// Tree godoc
// @Summary Get product category tree
// @Tags Product categories
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[[]v1.Category] "Successful response"
// @Router /categories/tree [get]
func (h *CategoryHandler) Tree(c *gin.Context) {
	items, err := h.svc.Tree(c)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	out := make([]v1.Category, len(items))
	for i := range items {
		out[i] = categoryJSON(items[i])
	}
	v1.HandleSuccess(c, out)
}
