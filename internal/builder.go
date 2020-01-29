package internal

// Builder is responsible for building an artifact from source.
type Builder interface {
	Build() (Artifact, error)
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

// PackBuilder uses `pack`, the frontend for CloudNative Buildpacks,
// to build an artifact from source.
type PackBuilder struct{}

func (b *PackBuilder) Build() (Artifact, error) {
}
