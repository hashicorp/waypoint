package mapper

import (
	"fmt"
	"reflect"
)

// Func is a dependency-injected function.
//
// The concrete function must return 1 or 2 results: a value alone or
// a value with an error. No other shape of function is allowed.
type Func struct {
	Args []reflect.Type
	Out  reflect.Type
	Func reflect.Value
}

// NewFunc creates a new Func from f. f must be a function pointer.
func NewFunc(f interface{}) (*Func, error) {
	ft := reflect.TypeOf(f)
	if k := ft.Kind(); k != reflect.Func {
		return nil, fmt.Errorf("fn should be a function, got %s", k)
	}

	// We should have one or two results: a concrete type alone or a concrete
	// type along with an error value.
	if n := ft.NumOut(); n == 0 || n > 2 {
		return nil, fmt.Errorf("fn should return one or two results, got %d", n)
	}

	// Build our args list
	args := make([]reflect.Type, 0, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		args = append(args, ft.In(i))
	}

	return &Func{
		Func: reflect.ValueOf(f),
		Args: args,
		Out:  ft.Out(0),
	}, nil
}

// Call calls the function with the given values and returns the result and
// any error. If all arguments aren't satisfied with values, an error will
// be returned. You can check in advance if arguments will work by calling
// Prepare and subsequently calling Call on the prepared func.
func (f *Func) Call(values ...interface{}) (interface{}, error) {
	pf := f.Prepare(values...)
	if pf == nil {
		return nil, fmt.Errorf("failed to call function, unsatisfied arguments")
	}

	return pf.Call()
}

// Prepare verifies that all the values satisfy the function and returns
// a PreparedFunc that is ready to be called with Call. If the values do not
// satisfy the parameter requirements of the function, nil is returned.
//
// Prepare performs some preprocessing required to call Call, so it is more
// performant to call Call on the returned PreparedFunc rather than Call on
// Func. This will avoid duplicate work.
func (f *Func) Prepare(values ...interface{}) *PreparedFunc {
	return f.prepare(f.valueMap(values...))
}

// valueMap turns a list of values into the map required for processing.
func (f *Func) valueMap(values ...interface{}) map[reflect.Type]reflect.Value {
	vt := make(map[reflect.Type]reflect.Value)
	for _, value := range values {
		v := reflect.ValueOf(value)
		vt[v.Type()] = v
	}

	return vt
}

// args builds the list of args for the function given the valueMap. If
// missing is non-nil, then it will be populated with any missing values.
// If missing is nil, then nil will be returned if there are any missing values.
func (f *Func) args(
	vt map[reflect.Type]reflect.Value,
	missing map[reflect.Type]int,
) []reflect.Value {
	in := make([]reflect.Value, len(f.Args))
	for idx, arg := range f.Args {
		v := vt[arg]

		if !v.IsValid() {
			// If we didn't find a direct type matching, then we go loop
			// through all the values to see if we have a value that implements
			// the interface argument. We only do this for interface types.
			if arg.Kind() == reflect.Interface {
				for t, vv := range vt {
					if t.Implements(arg) {
						v = vv
						break
					}
				}
			}

			// We didn't find a direct value or a value impl the interface
			if !v.IsValid() {
				if missing == nil {
					return nil
				}

				missing[arg] = idx
			}
		}

		// Store the argument
		in[idx] = v
	}

	return in
}

// prepare is an unexported version of Prepare that takes a further
// preprocessed value map form.
func (f *Func) prepare(vt map[reflect.Type]reflect.Value) *PreparedFunc {
	in := f.args(vt, nil)
	return &PreparedFunc{Func: f, In: in}
}

// PreparedFunc is created by calling Prepare on a Func and is a preprocessed
// form of a function ready to be called with pre-known arguments.
type PreparedFunc struct {
	Func *Func
	In   []reflect.Value
}

// Call calls the function and returns the results, see Func.Call.
func (f *PreparedFunc) Call() (interface{}, error) {
	var err error
	out := f.Func.Func.Call(f.In)
	if len(out) > 1 {
		if v := out[1].Interface(); v != nil {
			err = v.(error)
		}
	}

	return out[0].Interface(), err
}
