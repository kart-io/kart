package main

import (
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
