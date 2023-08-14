package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kart-io/kart/cmd/kart/app"
	"github.com/spf13/cobra"

	"github.com/kart-io/kart/cmd/kart/app"
)

var targetDir string

func Command() *cobra.Command {
	command := app.NewCommand("run", "This is the run command", func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	})
	command.Flags().StringVarP(&targetDir, "work", "w", "", "target working directory")
	return command
}

// Run project.
func Run(cmd *cobra.Command, args []string) {
	var dir string
	cmdArgs, programArgs := splitArgs(cmd, args)
	if len(cmdArgs) > 0 {
		dir = cmdArgs[0]
	}
	base, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		return
	}
	if dir == "" {
		cmdPath, err := findCMD(base)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		switch len(cmdPath) {
		case 0:
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current directory")
		case 1:
			for _, v := range cmdPath {
				dir = v
			}
		default:
			cmdPaths := make([]string, len(cmdPath))
			for k := range cmdPath {
				cmdPaths = append(cmdPaths, k)
			}
			prompt := &survey.Select{
				Message:  "Which directory do you want to run?",
				Options:  cmdPaths,
				PageSize: 10,
			}
			e := survey.AskOne(prompt, &dir)
			if e != nil || dir == "" {
				return
			}
			dir = cmdPath[dir]
		}
	}
	fd := exec.Command("go", append([]string{"run", dir}, programArgs...)...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = dir
	changeWorkingDirectory(fd, targetDir)
	if err := fd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}
	return
}

func splitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}

func findCMD(base string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}
	var root bool
	next := func(dir string) (map[string]string, error) {
		cmdPath := make(map[string]string)
		err := filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
			// multi level directory is not allowed under the cmdPath directory, so it is judged that the path ends with cmdPath.
			if strings.HasSuffix(walkPath, "cmd") {
				paths, err := os.ReadDir(walkPath)
				if err != nil {
					return err
				}
				for _, fileInfo := range paths {
					if fileInfo.IsDir() {
						abs := filepath.Join(walkPath, fileInfo.Name())
						cmdPath[strings.TrimPrefix(abs, wd)] = abs
					}
				}
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return cmdPath, err
	}
	for i := 0; i < 5; i++ {
		tmp := base
		cmd, err := next(tmp)
		if err != nil {
			return nil, err
		}
		if len(cmd) > 0 {
			return cmd, nil
		}
		if root {
			break
		}
		_ = filepath.Join(base, "..")
	}
	return map[string]string{"": base}, nil
}

func changeWorkingDirectory(cmd *exec.Cmd, targetDir string) {
	targetDir = strings.TrimSpace(targetDir)
	if targetDir != "" {
		cmd.Dir = targetDir
	}
}
