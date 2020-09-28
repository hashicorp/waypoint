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

//go:generate stringer -type=Type -linecomment
//go:generate mockery -all -case underscore

// Type is an enum of all the types of components supported.
// This isn't used directly in this package but is used by other packages
// to reference the component types.
type Type uint

const (
	InvalidType        Type = iota // Invalid
	BuilderType                    // Builder
	RegistryType                   // Registry
	PlatformType                   // Platform
	ReleaseManagerType             // ReleaseManager
	LogPlatformType                // LogPlatform
	LogViewerType                  // LogViewer
	AuthenticatorType              // Authenticator
	MapperType                     // Mapper
	maxType
)

// TypeMap is a mapping of Type to the nil pointer to the interface of that
// type. This can be used with libraries such as mapper.
var TypeMap = map[Type]interface{}{
	BuilderType:        (*Builder)(nil),
	RegistryType:       (*Registry)(nil),
	PlatformType:       (*Platform)(nil),
	ReleaseManagerType: (*ReleaseManager)(nil),
	LogPlatformType:    (*LogPlatform)(nil),
	LogViewerType:      (*LogViewer)(nil),
	AuthenticatorType:  (*Authenticator)(nil),
}

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

// PlatformReleaser is an optional interface that a Platform can implement
// to provide default Release functionality. This only takes effect if
// no release is configured.
type PlatformReleaser interface {
	// DefaultReleaserFunc() should return a function that returns
	// a ReleaseManger implementation. This ReleaseManager will NOT have
	// any config so it must work by default.
	DefaultReleaserFunc() interface{}
}

// ReleaseManager is responsible for taking a deployment and making it
// "released" which means that traffic can now route to it.
type ReleaseManager interface {
	// ReleaseFunc should return the method handle for the "release" operation.
	ReleaseFunc() interface{}
}

// Destroyer is responsible for destroying resources associated with this
// implementation. This can be implemented by all of the component types
// and will be called to perform cleanup on any created resources.
type Destroyer interface {
	// DestroyFunc should return the method handle for the destroy operation.
	DestroyFunc() interface{}
}

// WorkspaceDestroyer is called when a workspace destroy operation is
// performed (typically via the "waypoint destroy" CLI). This can be implemented
// by any plugin.
type WorkspaceDestroyer interface {
	// DestroyWorkspaceFunc is called when a workspace destroy operation is performed.
	//
	// This will only be called if that plugin had performed some operation
	// previously on the workspace. This may be called multiple times so it should
	// be idempotent. This will be called after all individual DestroyFuncs are
	// complete.
	DestroyWorkspaceFunc() interface{}
}

// Authenticator is responsible for authenticating different types of plugins.
type Authenticator interface {
	// AuthFunc should return the method for getting credentials for a
	// plugin. This should return AuthResult.
	AuthFunc() interface{}

	// ValidateAuthFunc should return the method for validating authentication
	// credentials for the plugin
	ValidateAuthFunc() interface{}
}

// See Args.Source in the protobuf protocol.
type Source struct {
	App  string
	Path string
}

// AuthResult is the return value expected from Authenticator.AuthFunc.
type AuthResult struct {
	// Authenticated when true means that the plugin should now be authenticated
	// (given the other fields in this struct). If ValidateAuth is called,
	// it should succeed. If this is false, the auth method may have printed
	// help text or some other information, but it didn't authenticate. However,
	// this is not an error.
	Authenticated bool
}

type LabelSet struct {
	Labels map[string]string
}

// JobInfo is available to plugins to get information about the context
// in which a job is executing.
type JobInfo struct {
	// Id is the ID of the job that is executing this plugin operation.
	// If this is empty then it means that the execution is happening
	// outside of a job.
	Id string

	// Local is true if the operation is running locally on a machine
	// alongside the invocation. This can be used to determine if you can
	// do things such as open browser windows, read user files, etc.
	Local bool

	// Workspace is the workspace that this job is executing in. This should
	// be used by plugins to properly isolate resources from each other.
	Workspace string
}

type Artifact interface {
	// Labels are the labels to set. These will overwrite any conflicting
	// labels on the value. Please namespace the labels you set. The recommended
	// namespacing is using a URL structure, followed by a slash, and a key.
	// For example: "plugin.example.com/key" as the key. The value can be
	// any string.
	Labels() map[string]string
}

type Deployment interface{}

type Release interface {
	// URL is the URL to access this release.
	URL() string
}
