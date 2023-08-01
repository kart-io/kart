package config

import (
	"github.com/kart-io/kart/pkg/app"
)

type Config struct {
	*app.Options
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given IAM pump command line or configuration file option.
func CreateConfigFromOptions(opts *app.Options) (*Config, error) {
	return &Config{opts}, nil
}
