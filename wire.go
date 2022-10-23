//go:build wireinject

//go:generate wire .
package webshop

import (
	"github.com/google/wire"
	"webshop/internal/config"
	"webshop/internal/database"
	"webshop/shop"
	"webshop/shop/boundary"
)

type AppContext struct {
	boundary.Router
}

func Initialize() *AppContext {
	panic(wire.Build(
		config.ProvideConfig,
		database.ProvideDatabase,

		shop.Providers,

		wire.Struct(new(AppContext), "*"),
	))
}
