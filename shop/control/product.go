package control

import (
	"context"
	"log"
	"time"
	"webshop/shop/entity"
)

type (
	ProductRepository interface {
		// Save persist a product in the database and returns its ID
		Save(ctx context.Context, product entity.Product) error

		// FindPaginated finds a paginated list of products. It can also be sorted and filtered via an entity.FilterObject
		FindPaginated(ctx context.Context, filterObject entity.FilterObject) ([]entity.Product, error)
	}

	ProductController struct {
		repository ProductRepository
	}
)

func ProvideController(repository ProductRepository) ProductController {
	return ProductController{repository: repository}
}

// CreateProduct creates a new product and returns its id
func (controller ProductController) CreateProduct(ctx context.Context, product entity.Product) (string, error) {
	if err := product.Validate(); err != nil {
		log.Printf("failed to validate product, %s", err)
		return "", err
	}

	now := time.Now()
	product.CreatedAt = now
	product.ModifiedAt = now
	product.NewId()

	if err := controller.repository.Save(ctx, product); err != nil {
		log.Printf("failed to persist product with title '%s', %s", product.Title, err)
		return "", err
	}

	log.Printf("created product with title '%s'", product.Title)

	return product.ID, nil
}

func (controller ProductController) FindProducts(ctx context.Context, filter entity.FilterObject) ([]entity.Product, error) {
	if err := filter.Validate(); err != nil {
		log.Printf("failed to validate filter object, %s", err)
		return nil, err
	}

	products, err := controller.repository.FindPaginated(ctx, filter)
	if err != nil {
		log.Printf("failed to find products, %s", err)
		return nil, err
	}

	log.Printf("found '%d' products", len(products))

	return products, nil
}
