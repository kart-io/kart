package app

import (
	"fmt"
	"github.com/kart-io/kart/cmd/kart/version"
	"github.com/spf13/cobra"
)

type Option func(*App)

type App struct {
	basename    string
	description string
	cmd         *cobra.Command
}

func WithCommand(cmds ...*cobra.Command) Option {
	return func(a *App) {
		a.cmd.AddCommand(cmds...)
	}
}

func WithBasename(name string) Option {
	return func(a *App) {
		a.basename = name
	}
}

func WithDescription(description string) Option {
	return func(a *App) {
		a.description = description
	}
}

func NewApp(basename string, options ...Option) *App {
	app := &App{
		cmd: &cobra.Command{
			Use:     basename,
			Short:   fmt.Sprintf("%s is a tool for go microservices", basename),
			Version: version.Release,
		},
	}
	for _, option := range options {
		option(app)
	}
	return app
}

func (a *App) Execute() error {
	return a.cmd.Execute()
}
