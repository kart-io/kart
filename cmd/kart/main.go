package main

import (
	"github.com/kart-io/kart/cmd/kart/app"
	"github.com/kart-io/kart/cmd/kart/internal/run"
	"log"
)

func main() {
	newApp := app.NewApp(
		"kart",
		app.WithCommand(
			run.Command(),
		),
	)
	if err := newApp.Execute(); err != nil {
		log.Fatal(err)
	}
}
