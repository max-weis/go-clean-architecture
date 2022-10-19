package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"webshop/internal/config"
)

func ProvideDatabase(cfg config.Config) *sqlx.DB {
	dbCfg := cfg.DatabaseConfig
	connection := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DbName)

	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
