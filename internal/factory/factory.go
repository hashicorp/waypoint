// Package factory contains a "factory" pattern based on argmapper.
//
// A Factory can be used to register factory methods to create some predefined
// type or interface implementation. These functions are argmapper functions so
// converters and so on can be used as part of instantiation.
package factory

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
)

// Factory keeps track of named dependency-injected factory functions to
// create an implementation of an interface.
type Factory struct {
	iface reflect.Type
	funcs map[string]*argmapper.Func
}

// NewFactory creates a Factory for the interface iface. The parameter
// iface should be a nil pointer to the interface type. Example: (*iface)(nil).
func NewFactory(iface interface{}) (*Factory, error) {
	// Get the interface type
	it := reflect.TypeOf(iface)
	if k := it.Kind(); k != reflect.Ptr {
		return nil, fmt.Errorf("iface must be a pointer to an interface, got %s", k)
	}
	it = it.Elem()
	if k := it.Kind(); k != reflect.Interface {
		return nil, fmt.Errorf("iface must be a pointer to an interface, got %s", k)
	}

	return &Factory{iface: it, funcs: make(map[string]*argmapper.Func)}, nil
}

// Register registers a factory function named name for the interface.
func (f *Factory) Register(name string, fn interface{}) error {
	ff, err := argmapper.NewFunc(fn)
	if err != nil {
		return err
	}

	f.funcs[name] = ff
	return nil
}

// Func returns the factory function named name. This can then be used to
// call and instantiate the factory interface type.
func (f *Factory) Func(name string) *argmapper.Func {
	return f.funcs[name]
}
