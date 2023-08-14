package command

import (
	"os"

	"github.com/costa92/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kart-io/kart/pkg/cliflag"
)

var progressMessage = color.GreenString("==>")

type App struct {
	basename    string
	options     CliOptions
	runFunc     RunFunc
	description string
	commands    []*Command
	cmd         *cobra.Command
	noConfig    bool
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// WithOptions returns a function that calls each Option in Options
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithRunFunc returns a function that calls each RunFunc in RunFuncs
func WithRunFunc(runFunc RunFunc) Option {
	return func(a *App) {
		a.runFunc = runFunc
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

func NewApp(basename string, opts ...Option) *App {
	a := &App{
		basename: basename,
	}
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
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, c := range a.commands {
			cmd.AddCommand(c.cobraCommand())
		}
		cmd.SetHelpCommand(HelpCommand(FormatBaseName(a.basename)))
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
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}
	a.cmd = &cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	PrintFlags(cmd.Flags())
	// 获取配置文件
	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}
	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}
	return nil
}

// Run is used to launch the application.
func (a *App) Run() error {
	if err := a.cmd.Execute(); err != nil {
		logger.Infof("%v %v\n", color.RedString("Error:"), err)
		return err
	}
	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	logger.Infof("%v WorkingDir: %s", progressMessage, wd)
}
