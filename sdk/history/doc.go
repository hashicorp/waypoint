// Package history provides an interface for components to query historical
// actions. An example use case is for a platorm to look up historical
// deployments it may use determine whether it is creating or updating.
package history

//go:generate mockery -all -case underscore
