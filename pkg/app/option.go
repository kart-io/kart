package app

import "github.com/kart-io/kart/pkg/cliflag"

type Options struct{}

func NewOptions() *Options {
	o := Options{}
	return &o
}

func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	return fss
}
