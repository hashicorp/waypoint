package history

import (
	"context"

	"github.com/mitchellh/devflow/sdk/component"
)

// Client is the client to access historical information. Component
// plugins can add this as an argument to get access to an implementation.
type Client interface {
	// Deployments looks up past deployments.
	Deployments(context.Context, *Lookup) ([]component.Deployment, error)
}

// Lookup is the lookup configuration used by the history client.
type Lookup struct {
	// Type should be a concrete type that all results should be converted
	// to. Example value: (*Foo)(nil) where Foo is a struct. This will cause
	// an error if any result can't map to the struct given. If this is set
	// and there is a result, the result is guaranteed to be this type.
	Type interface{}

	// Limit is the number of results to return
	Limit int

	// FilterStatus allows filtering by as specific status.
	FilterStatus FilterStatus
}

// FilterStatus is a value for what status to look for in a lookup.
type FilterStatus uint

const (
	StatusInvalid FilterStatus = iota // invalid
	StatusSuccess                     // success
	StatusError                       // error
)
