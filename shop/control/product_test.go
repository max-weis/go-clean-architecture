package control

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		assert.ErrorIs(t, err, entity.ErrValidation)
	})

	t.Run("Invalid description", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		_, err := controller.CreateProduct(context.Background(), entity.Product{
			Title:       "test",
			Description: "",
			Price:       10,
		})
		assert.ErrorIs(t, err, entity.ErrValidation)
	})
}

func TestProductController_FindProducts(t *testing.T) {
	t.Run("Successfully find 10 products", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		products, err := controller.FindProducts(context.Background(), entity.FilterObject{
			Sort:   entity.None,
			Limit:  10,
			Offset: 0,
		})
		assert.NoError(t, err)

		assert.Len(t, products, 10)
	})

	t.Run("Invalid limit", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		_, err := controller.FindProducts(context.Background(), entity.FilterObject{
			Sort:   entity.None,
			Limit:  0,
			Offset: 0,
		})
		assert.ErrorIs(t, err, entity.ErrValidation)
	})
}

func TestProductController_FindProduct(t *testing.T) {
	t.Run("Successfully find product", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		product, err := controller.FindProduct(context.Background(), "test")
		assert.NoError(t, err)

		assert.Equal(t, "test", product.ID)
		assert.Equal(t, "test", product.Title)
		assert.Equal(t, "test", product.Description)
	})

	t.Run("Product not found", func(t *testing.T) {
		controller := ProvideController(mockProductRepository{})

		_, err := controller.FindProduct(context.Background(), "")
		assert.ErrorIs(t, err, entity.ErrProductNotFound)
	})
}

type mockProductRepository struct {
}

func (m mockProductRepository) FindByID(ctx context.Context, id string) (entity.Product, error) {
	if id == "" {
		return entity.Product{}, nil
	}

	return entity.Product{
		ID:          "test",
		Title:       "test",
		Description: "test",
		Price:       1,
		CreatedAt:   time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC),
		ModifiedAt:  time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC),
	}, nil
}

func (m mockProductRepository) FindPaginated(ctx context.Context, filterObject entity.FilterObject) ([]entity.Product, error) {
	return []entity.Product{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}, {ID: "6"}, {ID: "7"}, {ID: "8"}, {ID: "9"}, {ID: "10"}}, nil
}

func (m mockProductRepository) Save(ctx context.Context, product entity.Product) error {
	return nil
}
