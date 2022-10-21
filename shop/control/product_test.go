package control

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"webshop/shop/entity"
)

func TestProductController_CreateProduct(t *testing.T) {
	t.Run("Successfully created product", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		id, err := controller.CreateProduct(context.Background(), entity.Product{
			Title:       "test",
			Description: "test",
			Price:       10,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("Invalid title", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		_, err := controller.CreateProduct(context.Background(), entity.Product{
			Title:       "",
			Description: "test",
			Price:       10,
		})
		assert.ErrorIs(t, err, ValidationError)
	})

	t.Run("Invalid description", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		_, err := controller.CreateProduct(context.Background(), entity.Product{
			Title:       "test",
			Description: "",
			Price:       10,
		})
		assert.ErrorIs(t, err, ValidationError)
	})
}

type mockProductRepository struct {
}

func (m mockProductRepository) Save(ctx context.Context, product entity.Product) error {
	return nil
}
