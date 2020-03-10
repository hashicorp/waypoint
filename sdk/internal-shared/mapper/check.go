package mapper

import (
	"reflect"
)

// CheckFunc is a function type used with many Chain functions to perform
// a check if a type is satisfied. The exact meaning of "satisfied" depends
// on the chain function called.
type CheckFunc func(Type) bool

// CheckEqual returns a CheckFunc that compares that two types are equal.
// Equal is defined by their Key value being the same.
func CheckEqual(typ Type) CheckFunc {
	return func(cmp Type) bool { return cmp.Key() == typ.Key() }
}

// CheckReflectType returns a CheckFunc that returns true if the type
// matches the Go reflect type. If the given type is an interface, then
// it will also return true if the type implements that interface.
func CheckReflectType(t reflect.Type) CheckFunc {
	return func(cmp Type) bool {
		rt, ok := cmp.(*ReflectType)
		if !ok {
			return false
		}

		if rt.Type == t {
			return true
		}

		if t.Kind() == reflect.Interface && rt.Type.Implements(t) {
			return true
		}

		return false
	}
}

// CheckOr composes any number of CheckFuncs with a logical OR.
func CheckOr(fs ...CheckFunc) CheckFunc {
	return func(t Type) bool {
		for _, f := range fs {
			if f(t) {
				return true
			}
		}

		return false
	}
}
