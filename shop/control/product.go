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

		// Update updates a product with the given id
		Update(ctx context.Context, product entity.Product) error

		// Delete removes a product with the given id
		Delete(ctx context.Context, id string) error

		// FindPaginated finds a paginated list of products. It can also be sorted and filtered via an entity.FilterObject
		FindPaginated(ctx context.Context, filterObject entity.FilterObject) ([]entity.Product, error)

		// FindByID tries to find a product.
		// Returns an entity.ErrProductNotFound error if no product exists with the given id
		FindByID(ctx context.Context, id string) (entity.Product, error)
	}

	ProductController struct {
		repository ProductRepository
	}
)

func ProvideController(repository ProductRepository) ProductController {
	return ProductController{repository: repository}
}

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

func (controller ProductController) UpdateProduct(ctx context.Context, id string, product entity.Product) error {
	update, err := controller.FindProduct(ctx, id)
	if err != nil {
		return err
	}

	update.Title = product.Title
	update.Description = product.Description
	update.Price = product.Price
	update.ModifiedAt = time.Now()

	if err := controller.repository.Update(ctx, update); err != nil {
		return err
	}

	log.Printf("updated product with id '%s'", product.ID)

	return nil
}

func (controller ProductController) DeleteProduct(ctx context.Context, id string) error {
	_, err := controller.FindProduct(ctx, id)
	if err != nil {
		return err
	}

	if err := controller.repository.Delete(ctx, id); err != nil {
		return err
	}

	log.Printf("deleted product with id '%s'", id)

	return nil
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

func (controller ProductController) FindProduct(ctx context.Context, id string) (entity.Product, error) {
	product, err := controller.repository.FindByID(ctx, id)
	if err != nil {
		log.Printf("failed to find product with id '%s', %s", id, err)
		return entity.Product{}, err
	}

	if product.ID == "" {
		return entity.Product{}, entity.ErrProductNotFound
	}

	log.Printf("found product with id '%s'", id)

	return product, nil
}
