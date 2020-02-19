package component

import "errors"

type ConfigVar struct {
	Name, Value string
}

var ErrNoSuchVariable = errors.New("no such variable exists")
