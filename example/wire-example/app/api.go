package app

import (
	"github.com/google/wire"
	"github.com/kart-io/kart/example/wire-example/config"
	"github.com/kart-io/kart/example/wire-example/options"
	"github.com/kart-io/kart/internal/command"
)

const commandDesc = `The Kart API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.

Find more api server information at:
    https://github.com/costa92/kart`

var ProviderApiServerSet = wire.NewSet(NewOptionConfig, NewApiServer)

type ApiServer struct {
	App *command.App
}

func NewApiServer(opts *options.Options) *ApiServer {
	a := command.NewApp(
		"kart",
		command.WithOptions(opts),
		command.WithDescription(commandDesc),
		command.WithRunFunc(run(opts)),
	)
	return &ApiServer{
		App: a,
	}
}

func (a *ApiServer) Run() error {
	return a.App.Run()
}

func NewOptionConfig() *options.Options {
	return options.NewOptions()
}

func run(opts *options.Options) command.RunFunc {
	return func(basename string) error {
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return nil
		}
		return Run(cfg)
	}
}
