/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"

	"github.com/kart-io/kart/internal/command"
	"github.com/kart-io/kart/internal/config"
	"github.com/kart-io/kart/pkg/app"
)

type Config struct {
	Server Server `json:"server,omitempty" yaml:"server" mapstructure:"server" toml:"server"`
}

type Server struct {
	Port string `json:"port,omitempty" yaml:"port" mapstructure:"port" toml:"port"`
}

const commandDesc = `The Kart API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.

Find more iam-apiserver information at:
    https://github.com/costa92/kart`

func main() {
	opts := app.NewOptions()
	a := command.NewApp(
		"kart",
		command.WithOptions(opts),
		command.WithDescription(commandDesc),
		command.WithRunFunc(run(opts)),
	)
	err := a.Run()
	if err != nil {
		return
	}
}

func run(opts *app.Options) command.RunFunc {
	return func(basename string) error {
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		return Run(cfg)
	}
}

func Run(config *config.Config) error {
	fmt.Println(config)
	return nil
}
