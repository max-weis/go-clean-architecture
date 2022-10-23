package boundary

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"strings"
	"testing"
	"time"
	"webshop/internal/config"
	"webshop/internal/database"
	"webshop/shop/entity"
)

func TestSqlxRepository_Query(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrationtest")
	}

	compose := setupPostgres(t)
	defer compose.Down()

	// TODO: use waiting strategy
	time.Sleep(5 * time.Second)

	repository := setupRepository()

	t.Run("Find products", func(t *testing.T) {
		products, err := repository.FindPaginated(context.Background(), entity.FilterObject{
			Sort:   entity.TitleAsc,
			Limit:  10,
			Offset: 0,
		})
		assert.NoError(t, err)

		assert.Len(t, products, 10)
	})

	t.Run("Find product", func(t *testing.T) {
		product, err := repository.FindByID(context.Background(), "564a495c-3ac1-475d-b257-bd022fad7f96")
		assert.NoError(t, err)

		assert.Equal(t, product.ID, "564a495c-3ac1-475d-b257-bd022fad7f96")
		assert.Equal(t, product.Title, "Pasta - Lasagne, Fresh")
	})

	t.Run("Found no product", func(t *testing.T) {
		product, err := repository.FindByID(context.Background(), "1")
		assert.NoError(t, err)

		assert.Equal(t, product, entity.Product{})
	})
}

func TestSqlxRepository_Command(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrationtest")
	}

	compose := setupPostgres(t)
	defer compose.Down()

	// TODO: use waiting strategy
	time.Sleep(5 * time.Second)

	repository := setupRepository()

	t.Run("Successfully save product", func(t *testing.T) {
		now := time.Now()
		err := repository.Save(context.Background(), entity.Product{
			ID:          "test",
			Title:       "test",
			Description: "test",
			Price:       10,
			CreatedAt:   now,
			ModifiedAt:  now,
		})
		assert.NoError(t, err)
	})

	t.Run("Successfully update", func(t *testing.T) {
		now := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
		ctx := context.Background()
		err := repository.Save(ctx, entity.Product{
			ID:          "test update",
			Title:       "test",
			Description: "test",
			Price:       10,
			CreatedAt:   now,
			ModifiedAt:  now,
		})
		assert.NoError(t, err)

		updateTime := time.Date(2000, time.January, 1, 22, 0, 0, 0, time.UTC)
		err = repository.Update(ctx, entity.Product{
			ID:          "test update",
			Title:       "test update",
			Description: "test update",
			Price:       99,
			CreatedAt:   updateTime,
			ModifiedAt:  updateTime,
		})

		product, err := repository.FindByID(ctx, "test update")
		assert.NoError(t, err)

		assert.Equal(t, product.ID, "test update")
		assert.Equal(t, product.Title, "test update")
		assert.Equal(t, product.Description, "test update")
		assert.Equal(t, product.Price, uint64(99))
	})

	t.Run("Successfully delete a product", func(t *testing.T) {
		now := time.Now()
		ctx := context.Background()
		err := repository.Save(ctx, entity.Product{
			ID:          "test delete",
			Title:       "test",
			Description: "test",
			Price:       10,
			CreatedAt:   now,
			ModifiedAt:  now,
		})
		assert.NoError(t, err)

		err = repository.Delete(ctx, "test delete")
		assert.NoError(t, err)

		product, err := repository.FindByID(ctx, "test update")
		assert.NoError(t, err)

		assert.Equal(t, product, entity.Product{})
	})
}

func setupRepository() sqlxRepository {
	return sqlxRepository{
		db: database.ProvideDatabase(config.ProvideConfig()),
	}
}

func setupPostgres(t *testing.T) testcontainers.DockerCompose {
	composeFilePaths := []string{"../../build/docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String())

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, identifier).WithCommand([]string{"up", "-d", "postgres"})
	if err := compose.Invoke().Error; err != nil {
		t.Fatal(err)
	}

	return compose
}
