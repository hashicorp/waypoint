package mapper

// CheckFunc is a function type used with many Chain functions to perform
// a check if a type is satisfied. The exact meaning of "satisfied" depends
// on the chain function called.
type CheckFunc func(Type) bool

// CheckEqual returns a CheckFunc that compares that two types are equal.
// Equal is defined by their Key value being the same.
func CheckEqual(typ Type) CheckFunc {
	return func(cmp Type) bool { return cmp.Key() == typ.Key() }
}
