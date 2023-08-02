package main

import (
	"github.com/kart-io/kart/cmd/kart/internal/options"
	"log"

	"github.com/kart-io/kart/cmd/kart/internal/run"
	"github.com/kart-io/kart/cmd/kart/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kart",
	Short:   "Kart is a tool for go microservices",
	Long:    "Kart is a tool for go microservices",
	Version: version.Release,
}

func init() {
	rootCmd.AddCommand(run.CmdRun)
}

// main is the entry point for the CLI

func main() {
	app := options.NewApp(
		"kart",
		options.WithDescription("Kart is a tool for go microservices"),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
