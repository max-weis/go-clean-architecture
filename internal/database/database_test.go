package database

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"strings"
	"testing"
	"time"
	"webshop/internal/config"
)

func TestProvideDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrationtest")
	}

	compose := setupPostgres(t)
	defer compose.Down()

	// TODO: use waiting strategy
	time.Sleep(5 * time.Second)

	cfg := config.ProvideConfig()

	db := ProvideDatabase(cfg)
	assert.NotNil(t, db)
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
