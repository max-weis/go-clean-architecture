package control

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
	"webshop/shop/entity"
)

var ValidationError = errors.New("validation failed")

type (
	ProductRepository interface {
		// Save persist a product in the database and returns its ID
		Save(ctx context.Context, product entity.Product) error
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
	if err := validateProduct(product); err != nil {
		log.Printf("failed to validate product, %s", err)
		return "", err
	}

	now := time.Now()
	product.CreatedAt = now
	product.ModifiedAt = now

	product.ID = uuid.NewString()

	if err := controller.repository.Save(ctx, product); err != nil {
		log.Printf("failed to persist product with title '%s', %s", product.Title, err)
		return "", err
	}

	log.Printf("created product with title '%s'", product.Title)

	return product.ID, nil
}

func validateProduct(product entity.Product) error {
	if product.Title == "" {
		return fmt.Errorf("%w: title must be set", ValidationError)
	}

	if product.Description == "" {
		return fmt.Errorf("%w: description must be set", ValidationError)
	}

	return nil
}
