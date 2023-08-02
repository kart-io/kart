package options

import (
	"github.com/costa92/logger"
	"github.com/fatih/color"
	"github.com/kart-io/kart/cmd/kart/internal/commands"
	"github.com/kart-io/kart/cmd/kart/version"
	"github.com/kart-io/kart/pkg/cliflag"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime"
	"strings"
)

var progressMessage = color.GreenString("==>")

type App struct {
	basename    string
	cmd         *cobra.Command
	options     commands.CliOptions
	commands    []*commands.Command
	description string
	runFunc     RunFunc
}

func NewApp(basename string, opts ...Option) *App {
	a := &App{}
	for _, o := range opts {
		o(a)
	}
	a.buildCommand()
	return a
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:           FormatBaseName(a.basename),
		Short:         "kart",
		Long:          a.description,
		Version:       version.Release,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	commands.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, c := range a.commands {
			cmd.AddCommand(c.CobraCommand())
		}
		cmd.SetHelpCommand(commands.HelpCommand(FormatBaseName(a.basename)))
	}

	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}
	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	a.cmd = &cmd
}

func (a *App) Run() error {
	if err := a.cmd.Execute(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// FormatBaseName is formatted as an executable file name under different
// operating systems according to the given name.
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}
	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	logger.Infof("%v WorkingDir: %s", progressMessage, wd)
}
