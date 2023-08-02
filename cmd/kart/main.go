package main

import (
	"github.com/kart-io/kart/cmd/kart/internal/options"
	"github.com/kart-io/kart/cmd/kart/internal/run"
	"log"
)

// main is the entry point for the CLI
func main() {
	basename := "kart"
	app := options.NewApp(
		basename,
		options.WithCommands(run.AddRun()),
		options.WithDescription("Kart is a tool for go microservices"),
		options.WithRunFunc(Run(basename)),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func Run(basename string) options.RunFunc {
	return func(basename string) error {
		return nil
	}
}
