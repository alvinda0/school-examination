package handlers

import (
	"net/http"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	meta := &api.Metadata{
		Page:       1,
		Limit:      10,
		Total:      len(products),
		TotalPages: 1,
	}
	utils.JSONResponse(w, http.StatusOK, "Products retrieved successfully", products, meta)
}
