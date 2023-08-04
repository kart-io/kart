package app

import "github.com/spf13/cobra"

type runFunc = func(cmd *cobra.Command, args []string)

func NewCommand(use string, short string, run runFunc) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run:   run,
	}
}
