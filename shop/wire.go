package shop

import (
	"github.com/google/wire"
	"webshop/shop/boundary"
	"webshop/shop/control"
)

var Providers = wire.NewSet(
	boundary.ProvideSqlxRepository,
	control.ProvideController,
	boundary.ProvideRouter,
)
