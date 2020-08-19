package config

import (
	"github.com/mitchellh/mapstructure"
)

// Operation is something in the Waypoint configuraiton that is executed
// using some underlying plugin. This is a general shared structure that is
// used by internal/core to initialize all the proper plugins.
type Operation struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

func (b *Build) Operation() *Operation {
	return mapoperation(b)
}

func (b *Build) RegistryOperation() *Operation {
	if b == nil {
		return nil
	}

	return b.Registry.Operation()
}

func (b *Registry) Operation() *Operation {
	return mapoperation(b)
}

func (b *Deploy) Operation() *Operation {
	return mapoperation(b)
}

func (b *Release) Operation() *Operation {
	return mapoperation(b)
}

// mapoperation takes a struct that is a superset of Operation and
// maps it down to an Operation. This will panic if this fails.
func mapoperation(input interface{}) *Operation {
	if input == nil {
		return nil
	}

	var op Operation
	if err := mapstructure.Decode(input, &op); err != nil {
		panic(err)
	}

	return &op
}
