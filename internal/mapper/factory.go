package mapper

import (
	"fmt"
	"reflect"
)

// Factory keeps track of named dependency-injected factory functions to
// create an implementation of an interface.
type Factory struct {
	iface reflect.Type
	funcs map[string]*Func
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

	return &Factory{iface: it, funcs: make(map[string]*Func)}, nil
}

// Register registers a factory function named name for the interface.
func (f *Factory) Register(name string, fn interface{}) error {
	ff, err := NewFunc(fn)
	if err != nil {
		return err
	}

	if !ff.Out.Implements(f.iface) {
		return fmt.Errorf("result of factory must implement interface")
	}

	f.funcs[name] = ff
	return nil
}

// Func returns the factory function named name. This can then be used to
// call and instantiate the factory interface type.
func (f *Factory) Func(name string) *Func {
	return f.funcs[name]
}
