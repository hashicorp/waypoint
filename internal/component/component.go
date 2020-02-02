// Package component has the interfaces for all the components that
// can be implemented. A component is the broad term used to describe
// all builders, platforms, registries, etc.
//
// Many component interfaces have functions named `XFunc` where "X" is some
// operation and the return value is "interface{}". These functions should return
// a method handle to the function implementing that operation. This pattern is
// done so that we can support custom typed operations that take and return
// full rich types for an operation. We use a minimal dependency-injection
// framework (see internal/mapper) to call these functions.
package component

import (
	"fmt"
)

//go:generate stringer -type=Type

// Type is an enum of all the types of components supported.
// This isn't used directly in this package but is used by other packages
// to reference the component types.
type Type uint

const (
	InvalidType Type = iota
	BuilderType
	RegistryType
	PlatformType
)

// Builder is responsible for building an artifact from source.
type Builder interface {
	// BuildFunc should return the method handle for the "build" operation.
	// The build function has access to a *Source and should return an Artifact.
	BuildFunc() interface{}
}

// Registry is responsible for managing artifacts.
type Registry interface {
	// PushFunc should return the method handle to the function for the "push"
	// operation. The push function should take an artifact type and push it
	// to the registry.
	PushFunc() interface{}
}

// Platform is responsible for deploying artifacts.
type Platform interface {
	// DeployFunc should return the method handle for the "deploy" operation.
	// The deploy function has access to the following and should use this
	// as necessary to perform a deploy.
	//
	//   artifact, artifact registry
	//
	DeployFunc() interface{}
}

// Source represents the source code for an application. This is used by
// the builder for creating an Artifact.
type Source struct {
	// App is the name of the application being built.
	App string

	// Path is the path to the root directory of the source tree.
	Path string
}

type Artifact interface{}

type Deployment interface {
	fmt.Stringer
}
