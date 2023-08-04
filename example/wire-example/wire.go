//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/kart-io/kart/example/wire-example/app"
)

func wireApp() (*app.ApiServer, func(), error) {
	panic(wire.Build(
		app.ProviderApiServerSet,
	))
}
