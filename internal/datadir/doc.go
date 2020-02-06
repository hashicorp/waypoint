// Package datadir manages the data directories. This includes persisted
// data such as state as well as ephemeral data such as cache and runtime
// files.
//
// This package is aware of the data model presented and provides easy
// helpers to create app-specific, component-specific, etc. data directories.
//
// This package is the result of lessons learned from reimplementing
// "data directories" for projects such as Vagrant and Terraform. Those
// projects managed a list of directories directly in the CLI, forcing
// a lot of code to be aware of paths and making it hard to implement
// operations on those paths such as pruning, migration, compression, etc.
// As an evolution, we create the "datadir" package which has deep knowledge
// of the software data model and consumers interact using higher level APIs
// rather than direct filesystem manipulation. This gives us more room to
// introduce improvements in the future that broadly impact the application
// without having to make those changes in many places.
package datadir
