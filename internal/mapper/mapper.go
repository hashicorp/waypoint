// Package mapper is able to perform basic dependency-injection-like
// functionality to support mapping from some set of source types to
// a desired destination type.
package mapper

import (
	"fmt"
	"reflect"
)

// M represents a single interface type and a set of mappers that can create
// concrete values of that interface type.
//
// The concrete implementations are stored with a human-friendly string name
// so that when attempting to create a target value, the desired target
// value type can be specified with a string.
type M struct {
	iface   reflect.Type
	impl    map[string]reflect.Type
	mappers []*mapper
}

// NewM creates an M for creating values that implement type iface.
//
// This will panic if an invalid value is given for iface. The value for iface
// should be a pointer to an interface with a nil value. Example:
//
//   (*myInterface)(nil)
//
func NewM(iface interface{}) *M {
	// Get the interface type
	it := reflect.TypeOf(iface)
	if k := it.Kind(); k != reflect.Ptr {
		panic(fmt.Sprintf("iface must be a pointer to an interface, got %s", k))
	}
	it = it.Elem()
	if k := it.Kind(); k != reflect.Interface {
		panic(fmt.Sprintf("iface must be a pointer to an interface, got %s", k))
	}

	return &M{iface: it, impl: make(map[string]reflect.Type)}
}

// RegisterImpl registers an implementation of the interface type tracked
// by this M. The value of impl should be a zero value of a struct implementing
// the interface. Example: (*myStruct)(nil)
func (m *M) RegisterImpl(name string, impl interface{}) error {
	// Get the concrete implementation
	ct := reflect.TypeOf(impl)
	if !ct.Implements(m.iface) {
		return fmt.Errorf("concrete (%s) must implement iface (%s)", ct, m.iface)
	}

	// Store it
	m.impl[name] = ct

	return nil

}

// RegisterMapper registers a mapper function to go from some argument
// types to a destination type.
func (m *M) RegisterMapper(name string, fn interface{}) error {
	ft := reflect.TypeOf(fn)
	if k := ft.Kind(); k != reflect.Func {
		return fmt.Errorf("fn should be a function, got %s", k)
	}

	// We should have one or two results: a concrete type alone or a concrete
	// type along with an error value.
	if n := ft.NumOut(); n == 0 || n > 2 {
		return fmt.Errorf("fn should return one or two results, got %d", n)
	}

	// Zeroth value should be our impl type and it should map the given name.
	ot := ft.Out(0)
	if m.impl[name] != ot {
		return fmt.Errorf("fn should output the same type expected for name %q", name)
	}

	// Build our mapper
	args := make([]reflect.Type, 0, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		args = append(args, ft.In(i))
	}

	m.mappers = append(m.mappers, &mapper{
		Args: args,
		Func: reflect.ValueOf(fn),
	})

	return nil
}

// Mapper returns a function that constructs an implementation of iface
// named name with the given available values.
func (m *M) Mapper(name string, values ...interface{}) func() (interface{}, error) {
	vt := make(map[reflect.Type]reflect.Value)
	for _, value := range values {
		v := reflect.ValueOf(value)
		vt[v.Type()] = v
	}

	for _, mapper := range m.mappers {
		pm := mapper.Prepare(vt)
		if pm != nil {
			return pm.Call
		}
	}

	return nil
}

// mapper keeps track of a single mapper function by storing some precomputed
// reflection values to make operation slightly more performant.
type mapper struct {
	Args []reflect.Type
	Func reflect.Value
}

// Prepare creates a preparedMapper that can be called. The input values
// vt are mapped to the expected arguments. If the input values do not satisfy
// the required parameters of the mapper, nil is returned.
func (m *mapper) Prepare(vt map[reflect.Type]reflect.Value) *preparedMapper {
	in := make([]reflect.Value, len(m.Args))
	for idx, arg := range m.Args {
		v := vt[arg]
		if !v.IsValid() {
			return nil
		}

		in[idx] = v
	}

	return &preparedMapper{Mapper: m, In: in}
}

// preparedMapper is a mapper that has the arguments already slotted and
// ready to be called.
type preparedMapper struct {
	Mapper *mapper
	In     []reflect.Value
}

// Call calls the mapper and returns the resulting created value and any error.
func (m *preparedMapper) Call() (interface{}, error) {
	var err error
	out := m.Mapper.Func.Call(m.In)
	if len(out) > 1 {
		err = out[1].Interface().(error)
	}

	return out[0].Interface(), err
}
