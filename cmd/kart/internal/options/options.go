package options

import "github.com/kart-io/kart/cmd/kart/internal/commands"

type Option func(*App)

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithOptions returns a function that calls each Option in Options
func WithOptions(opt commands.CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// WithRunFunc returns a function that calls each RunFunc in RunFuncs
func WithRunFunc(runFunc RunFunc) Option {
	return func(a *App) {
		a.runFunc = runFunc
	}
}

func WithDescription(description string) Option {
	return func(a *App) {
		a.description = description
	}
}
