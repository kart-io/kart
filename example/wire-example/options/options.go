package options

import (
	"github.com/kart-io/kart/pkg/cliflag"
	"github.com/kart-io/kart/transports/options"
)

type Options struct {
	ServerRunOption        *options.ServerRunOption        `json:"server"  mapstructure:"server"`
	InsecureServingOptions *options.InsecureServingOptions `json:"insecure"  mapstructure:"insecure"`
	FeatureOptions         *options.FeatureOptions         `json:"feature"  mapstructure:"feature"`
}

func NewOptions() *Options {
	o := Options{
		ServerRunOption:        options.NewServerRunOption(),
		InsecureServingOptions: options.NewInsecureServingOptions(),
		FeatureOptions:         options.NewFeatureOptions(),
	}
	return &o
}

func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	return fss
}
