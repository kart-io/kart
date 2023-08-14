package options

import (
	"github.com/gin-gonic/gin"

	kartHttp "github.com/kart-io/kart/transports/kart-http"
)

type ServerRunOption struct {
	Name        string   `json:"name" mapstructure:"name"`
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewServerRunOption() *ServerRunOption {
	return &ServerRunOption{
		Name:        "kart",
		Mode:        gin.ReleaseMode,
		Healthz:     true,
		Middlewares: []string{},
	}
}

func (s *ServerRunOption) ApplyTo(c *kartHttp.ServerConfig) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares

	return nil
}
