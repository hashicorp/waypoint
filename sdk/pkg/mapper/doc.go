// Package mapper is a minimal DI (dependency-injection) framework.
//
// mapper maps values to function parameters based only on their type. If
// the type matches, then it is injected. Names of parameters are not significant
// and it is not possible to map two values of the same type. This is due to
// the requirements we have for mapper not requiring this.
//
// The core of mapper is the Func struct, which maps a single function.
// Higher level structs like Factory then build on top of this to provide
// more functionality.
package mapper
