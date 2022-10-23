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

func TestSqlxRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrationtest")
	}

	t.Run("Successfully save product", func(t *testing.T) {
		compose := setupPostgres(t)
		defer compose.Down()

		// TODO: use waiting strategy
		time.Sleep(5 * time.Second)

		repository := setupRepository()

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
}

func TestSqlxRepository_FindPaginated(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrationtest")
	}

	t.Run("Find products", func(t *testing.T) {
		compose := setupPostgres(t)
		defer compose.Down()

		// TODO: use waiting strategy
		time.Sleep(5 * time.Second)

		repository := setupRepository()

		products, err := repository.FindPaginated(context.Background(), entity.FilterObject{
			Sort:   entity.TitleAsc,
			Limit:  10,
			Offset: 0,
		})
		assert.NoError(t, err)

		assert.Len(t, products, 10)
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
