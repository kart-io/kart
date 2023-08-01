package command

import "github.com/kart-io/kart/pkg/cliflag"

type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
}
