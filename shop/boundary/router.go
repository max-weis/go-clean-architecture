//go:generate oapi-codegen --config=oapi-codegen-config.yml spec.yml
package boundary

import (
	"encoding/json"
	"errors"
	"fmt"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
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

func mapToList(products []entity.Product, filterObject entity.FilterObject) ProductList {
	var list ProductList

	list.Limit = int(filterObject.Limit)
	list.Offset = int(filterObject.Offset)
	list.Products = len(products)

	var format string
	if filterObject.Free {
		format = "/v1/product?limit=%d&offset=%d&sort=%s&free=true"
	} else {
		format = "/v1/product?limit=%d&offset=%d&sort=%s"
	}

	list.Curr = fmt.Sprintf(format, filterObject.Limit, filterObject.Offset, mapEntityFilterToParam(filterObject.Sort))

	nextOffset := filterObject.Offset + 1
	list.Next = fmt.Sprintf(format, filterObject.Limit, nextOffset, mapEntityFilterToParam(filterObject.Sort))

	if filterObject.Offset == 0 {
		list.Prev = nil
	} else {
		prevOffset := filterObject.Offset - 1
		prev := fmt.Sprintf(format, filterObject.Limit, prevOffset, mapEntityFilterToParam(filterObject.Sort))
		list.Prev = &prev
	}

	data := make([]ProductListItem, 0, len(products))
	for _, product := range products {
		listItem := ProductListItem{
			Title:      product.Title,
			Price:      int(product.Price),
			CreatedAt:  openapi_types.Date{Time: product.CreatedAt},
			ModifiedAt: openapi_types.Date{Time: product.ModifiedAt},
		}

		data = append(data, listItem)
	}

	list.Data = data

	return list
}

func mapEntityFilterToParam(sort entity.Sorting) FindProductsParamsSort {
	switch sort {
	case entity.None:
		return None
	case entity.IdAsc:
		return IdAsc
	case entity.IdDesc:
		return IdDesc
	case entity.TitleAsc:
		return TitleAsc
	case entity.TitleDesc:
		return TitleDesc
	case entity.CreatedAtAsc:
		return CreatedAsc
	case entity.CreatedAtDesc:
		return CreatedDesc
	case entity.ModifiedAtAsc:
		return ModifiedAsc
	case entity.ModifiedAtDesc:
		return ModifiedDesc
	default:
		return None
	}
}

func (router Router) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product CreateProductJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := router.controller.CreateProduct(r.Context(), mapBodyToEntity(product))
	if err != nil {
		if errors.Is(err, entity.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	location := fmt.Sprintf("/v1/product/%s", id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func (router Router) UpdateProduct(w http.ResponseWriter, r *http.Request, id ProductID) {
	var product UpdateProductJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := router.controller.UpdateProduct(r.Context(), id, mapBodyToEntity(product)); err != nil {
		if errors.Is(err, entity.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	location := fmt.Sprintf("/v1/product/%s", id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusOK)
}

func (router Router) FindProduct(w http.ResponseWriter, r *http.Request, id string) {
	product, err := router.controller.FindProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrProductNotFound) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, mapToProduct(product))
}

func (router Router) FindProducts(w http.ResponseWriter, r *http.Request, params FindProductsParams) {
	filterObject := mapFilterParamsToEntity(params)
	products, err := router.controller.FindProducts(r.Context(), filterObject)
	if err != nil {
		if errors.Is(err, entity.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, mapToList(products, filterObject))
}

func mapToProduct(product entity.Product) Product {
	return Product{
		Id:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       int(product.Price),
		CreatedAt:   openapi_types.Date{Time: product.CreatedAt},
		ModifiedAt:  openapi_types.Date{Time: product.ModifiedAt},
	}
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf("failed to marshal error '%s'", err)
		return
	}
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
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

func mapFilterParamsToEntity(params FindProductsParams) entity.FilterObject {
	var filter entity.FilterObject

	filter.Limit = uint(params.Limit)
	filter.Offset = uint(params.Offset)

	filter.Sort = mapFilterParamsSortToEntity(params.Sort)

	var free bool
	if params.Free != nil {
		free = *params.Free
	}

	filter.Free = free

	return filter
}

func mapFilterParamsSortToEntity(param FindProductsParamsSort) entity.Sorting {
	switch param {
	case None:
		return entity.None
	case IdAsc:
		return entity.IdAsc
	case IdDesc:
		return entity.IdDesc
	case TitleAsc:
		return entity.TitleAsc
	case TitleDesc:
		return entity.TitleDesc
	case CreatedAsc:
		return entity.CreatedAtAsc
	case CreatedDesc:
		return entity.CreatedAtDesc
	case ModifiedAsc:
		return entity.ModifiedAtAsc
	case ModifiedDesc:
		return entity.ModifiedAtDesc
	default:
		return entity.None
	}
}
