package upgrade

import (
	"fmt"
	"github.com/kart-io/kart/cmd/kart/app"
	"github.com/kart-io/kart/cmd/kart/internal/base"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := app.NewCommand("run", "This is the run command", func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	})
	return command
}

// Run upgrade the kratos tools.
func Run(_ *cobra.Command, _ []string) {
	err := base.GoInstall(
		"github.com/kart-io/kart/cmd/kart@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
