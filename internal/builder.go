package internal

import (
	"context"
)

// Builder is responsible for building an artifact from source.
type Builder interface {
	Build(context.Context) (Artifact, error)
}

// An artifact is the result of a Builder and can be stored in an
// ArtifactRegistry and deployed to a Platform.
type Artifact interface{}

// Source represents the source code for an application. This is used by
// the builder for creating an Artifact.
type Source struct {
	// App is the name of the application being built.
	App string

	// Path is the path to the root directory of the source tree.
	Path string
}
