//go:generate oapi-codegen --config=oapi-codegen-config.yml spec.yml
package boundary

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"webshop/shop/control"
	"webshop/shop/entity"
)

type Router struct {
	controller control.ProductController
}

func ProvideRouter(controller control.ProductController) Router {
	return Router{controller: controller}
}

func (router Router) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product CreateProductJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := router.controller.CreateProduct(r.Context(), mapBodyToEntity(product))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	location := fmt.Sprintf("/v1/product/%s", id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf("failed to marshal error '%s'", err)
		return
	}
}

func mapBodyToEntity(product CreateProductJSONRequestBody) entity.Product {
	return entity.Product{
		Title:       product.Title,
		Description: product.Description,
		Price:       uint64(product.Price),
		CreatedAt:   time.Time{},
		ModifiedAt:  time.Time{},
	}
}
