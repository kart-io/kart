package upgrade

import (
	"fmt"
	"github.com/kart-io/kart/cmd/kart/app"
	"github.com/kart-io/kart/cmd/kart/internal/base"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := app.NewCommand("upgrade", "This is the upgrade command", func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	})
	return command
}

func Run(_ *cobra.Command, _ []string) {
	err := base.GoInstall(
		"github.com/kart-io/kart/cmd/kart@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
