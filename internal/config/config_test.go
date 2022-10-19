package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadEnvVarWithDefaults(t *testing.T) {
	assert.Equal(t, "postgres", readEnvVar("PQ_USER", "postgres"))
	assert.Equal(t, "postgres", readEnvVar("PQ_PASSWORD", "postgres"))
	assert.Equal(t, "postgres", readEnvVar("PQ_HOST", "postgres"))
	assert.Equal(t, "postgres", readEnvVar("PQ_PORT", "postgres"))
	assert.Equal(t, "postgres", readEnvVar("PQ_DB", "postgres"))
}

func TestReadEnvVarWithValues(t *testing.T) {
	t.Setenv("PQ_USER", "test")
	t.Setenv("PQ_PASSWORD", "test")
	t.Setenv("PQ_HOST", "test")
	t.Setenv("PQ_PORT", "test")
	t.Setenv("PQ_DB", "test")

	assert.Equal(t, "test", readEnvVar("PQ_USER", ""))
	assert.Equal(t, "test", readEnvVar("PQ_PASSWORD", ""))
	assert.Equal(t, "test", readEnvVar("PQ_HOST", ""))
	assert.Equal(t, "test", readEnvVar("PQ_PORT", ""))
	assert.Equal(t, "test", readEnvVar("PQ_DB", ""))
}

func TestNewDatabaseConfig(t *testing.T) {
	config := ProvideConfig()

	assert.Equal(t, "postgres", config.DatabaseConfig.User)
	assert.Equal(t, "postgres", config.DatabaseConfig.Password)
	assert.Equal(t, "localhost", config.DatabaseConfig.Host)
	assert.Equal(t, "5432", config.DatabaseConfig.Port)
	assert.Equal(t, "shop", config.DatabaseConfig.DbName)

	assert.Equal(t, "localhost:6379", config.CacheConfig.Address)
}
