package config

import (
	"os"
)

type (
	Config struct {
		DatabaseConfig DatabaseConfig
		CacheConfig    CacheConfig
	}

	DatabaseConfig struct {
		User     string
		Password string
		Host     string
		Port     string
		DbName   string
	}

	CacheConfig struct {
		Address string
	}
)

func ProvideConfig() Config {
	return Config{
		DatabaseConfig: newDatabaseConfig(),
		CacheConfig:    newCacheConfig(),
	}
}

func newDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		User:     readEnvVar("PQ_USER", "postgres"),
		Password: readEnvVar("PQ_PASSWORD", "password"),
		Host:     readEnvVar("PQ_HOST", "localhost"),
		Port:     readEnvVar("PQ_PORT", "5432"),
		DbName:   readEnvVar("PQ_DB", "postgres"),
	}
}

func newCacheConfig() CacheConfig {
	return CacheConfig{
		Address: readEnvVar("REDIS_ADDRESS", "localhost:6379"),
	}
}

func readEnvVar(key, def string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}

	return value
}
