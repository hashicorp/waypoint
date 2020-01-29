package pack

import (
	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal"
)

// NewBuilderFromSource creates a new Builder.
//
// This is a valid mapper that can be registered.
func NewBuilderFromSource(
	log hclog.Logger,
	source *internal.Source,
) *Builder {
	return &Builder{
		log:    log,
		source: source,
	}
}
