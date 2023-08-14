package main

import (
	"log"

	"github.com/kart-io/kart/cmd/kart/app"
	"github.com/kart-io/kart/cmd/kart/internal/image"
	"github.com/kart-io/kart/cmd/kart/internal/run"
	"github.com/kart-io/kart/cmd/kart/internal/upgrade"
)

func main() {
	newApp := app.NewApp(
		"kart",
		app.WithCommand(
			run.Command(),
			upgrade.Command(),
			image.Command(),
		),
	)
	if err := newApp.Execute(); err != nil {
		log.Fatal(err)
	}
}
