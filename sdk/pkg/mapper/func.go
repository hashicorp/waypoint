package mapper

import (
	"fmt"
	"reflect"
	"runtime"
)

// Func is a dependency-injected function.
//
// The concrete function must return 1 or 2 results: a value alone or
// a value with an error. No other shape of function is allowed.
type Func struct {
	Name string
	Args []Type
	Out  Type
	Func reflect.Value
}

// NewFunc creates a new Func from f. f must be a function pointer.
func NewFunc(f interface{}, opts ...Option) (*Func, error) {
	ft := reflect.TypeOf(f)
	if k := ft.Kind(); k != reflect.Func {
		return nil, fmt.Errorf("fn should be a function, got %s", k)
	}

	// Apply all the configs
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	// We should have one or two results: a concrete type alone or a concrete
	// type along with an error value.
	if n := ft.NumOut(); n == 0 || n > 2 {
		return nil, fmt.Errorf("fn should return one or two results, got %d", n)
	}

	// Get the reflect.Value for f since we'll need to store this but also
	// needs to get the address and so on.
	fv := reflect.ValueOf(f)

	// Try to get a name for the function
	var name string
	if rfunc := runtime.FuncForPC(fv.Pointer()); rfunc != nil {
		name = rfunc.Name()
	}

	// Build our args list
	args := make([]Type, 0, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		var arg Type
		typ := ft.In(i)

		// If we have a type matching override then use that.
		if cfg.WithType != nil {
			f, ok := cfg.WithType[typ]
			if ok {
				arg = f(i, typ)
			}
		}

		// Default to the reflection type
		if arg == nil {
			arg = &ReflectType{Type: typ}
		}

		args = append(args, arg)
	}

	return &Func{
		Name: name,
		Func: fv,
		Args: args,
		Out:  &ReflectType{Type: ft.Out(0)},
	}, nil
}

// config is the intermediary configuration used by Option to configure
// funcs during NewFunc.
type config struct {
	WithType map[reflect.Type]func(int, reflect.Type) Type
}

// Option is used to configure NewFunc
type Option func(*config)

// WithType replaces arguments of the given reflection type with the
// special Type implementation returned by the callback.
func WithType(typ reflect.Type, f func(int, reflect.Type) Type) Option {
	return func(c *config) {
		if c.WithType == nil {
			c.WithType = make(map[reflect.Type]func(int, reflect.Type) Type)
		}

		c.WithType[typ] = f
	}
}

// String returns the name of the function. If a name is not available, the
// raw type signature is returned instead.
func (f *Func) String() string {
	if f.Name != "" {
		return f.Name
	}

	return f.Func.String()
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
	// TODO: panic if missing values
	in := f.args(values, nil)
	return &PreparedFunc{Func: f, In: in}
}

// args builds the list of args for the function given the valueMap. If
// missing is non-nil, then it will be populated with any missing values.
// If missing is nil, then nil will be returned if there are any missing values.
func (f *Func) args(
	values []interface{},
	missing map[Type]int,
) []reflect.Value {
	in := make([]reflect.Value, len(f.Args))

	// NOTE(mitchellh): this is not very efficient at all, but I think in
	// general this won't be a hot-spot because we won't have many args
	// or values, and the cost will be amortized across the execution of
	// the program. If it does become an issue we can potentially optimize
	// through optional interface implementations by Value.
	for idx, arg := range f.Args {
		value := arg.Match(values...)
		if value != nil {
			in[idx] = reflect.ValueOf(value)
			continue
		}

		if missing != nil {
			missing[arg] = idx
		}
	}

	return in
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
