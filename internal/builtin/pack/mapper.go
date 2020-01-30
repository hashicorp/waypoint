package pack

import (
	"github.com/hashicorp/go-hclog"
)

// NewBuilder is the factory for the builder.
func NewBuilder(log hclog.Logger) *Builder {
	return &Builder{log: log}
}
