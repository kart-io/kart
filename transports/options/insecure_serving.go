package options

import (
	"net"
	"strconv"

	kartHttp "github.com/kart-io/kart/transports/kart-http"
)

type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
	}
}

func (s *InsecureServingOptions) ApplyTo(c *kartHttp.ServerConfig) error {
	c.InsecureServing = &kartHttp.InsecureServingInfo{
		Address: net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort)),
	}
	return nil
}
