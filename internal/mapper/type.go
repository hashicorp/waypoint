package mapper

import (
	"fmt"
	"reflect"
)

// Type can act as an argument or result of a mapper function.
//
// In most cases, callers don't need to be aware of this interface at all.
// Using this interface is an advanced feature that enables more dynamic
// functionality for the mapped functions. For example, a custom Type impl.
// can differentiate between different values of the same type. By default,
// mapper uses type-based (from reflection) values where if the type matches,
// the value matches.
type Type interface {
	// String value is for human-friendly error messages only.
	fmt.Stringer

	// Match should return a non-nil value if the input matches. The resulting
	// value will be used as the argument. In most cases, the result matches
	// the input when a match occurs.
	Match(values ...interface{}) interface{}

	// Key should return a unique comparable key for this type. T1.Key == T2.Key
	// when both T1 and T2 are "equal" Type values. This is used so that we can
	// do major performance improvements for chaining functions.
	Key() interface{}
}

// ReflectType implements Type based on a raw reflect.Type.
type ReflectType struct {
	Type reflect.Type
}

// Match implements Value by checking if in matches our reflect.Type
// in of two ways: a direct match or a concrete type that implements
// our interface if our type is an interface.
func (v *ReflectType) Match(values ...interface{}) interface{} {
	for _, value := range values {
		inval := reflect.ValueOf(value)
		intyp := inval.Type()

		// If we have a direct match, then use it as-is
		if intyp == v.Type {
			return value
		}

		// If we have an interface, we can attempt to see if the value implements
		// the interface. If our type is not an interface, then we can't check this.
		if v.Type.Kind() != reflect.Interface {
			continue
		}
		if intyp.Implements(v.Type) {
			return value
		}
	}

	return nil
}

// Key implements Type by returning the reflect.Type value directly since
// they are directly comparable.
func (v *ReflectType) Key() interface{} {
	return v.Type
}

// String implements fmt.Stringer
func (v *ReflectType) String() string {
	return v.Type.String()
}

var _ Type = (*ReflectType)(nil)
