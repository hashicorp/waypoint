package mapper

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/hashicorp/go-hclog"
)

// Func is a dependency-injected function.
//
// The concrete function must return 1 or 2 results: a value alone or
// a value with an error. No other shape of function is allowed.
type Func struct {
	Name   string
	Args   []Type
	Out    Type
	Func   reflect.Value
	Logger hclog.Logger
	Values []interface{} // extra values
}

// NewFuncList creates a slice of funcs from a slice of function pointers.
// The options if given are passed to each call to NewFunc for each function
// pointer.
func NewFuncList(fs []interface{}, opts ...Option) ([]*Func, error) {
	var result []*Func
	for _, raw := range fs {
		f, err := NewFunc(raw, opts...)
		if err != nil {
			return nil, err
		}

		result = append(result, f)
	}

	return result, nil
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

	// Defaults
	if cfg.Logger == nil {
		cfg.Logger = hclog.L()
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
	name := cfg.Name
	if name == "" {
		if rfunc := runtime.FuncForPC(fv.Pointer()); rfunc != nil {
			name = rfunc.Name()
		}
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
		Name:   name,
		Func:   fv,
		Args:   args,
		Out:    &ReflectType{Type: ft.Out(0)},
		Logger: cfg.Logger,
		Values: cfg.Values,
	}, nil
}

// config is the intermediary configuration used by Option to configure
// funcs during NewFunc.
type config struct {
	Name     string
	Logger   hclog.Logger
	WithType map[reflect.Type]func(int, reflect.Type) Type
	Values   []interface{}
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

// WithLogger sets a logger onto the mapper. If this isn't set it defaults to
// hclog.L() (the default logger).
func WithLogger(log hclog.Logger) Option {
	return func(c *config) { c.Logger = log }
}

// WithName sets the function name. This is used in debug messages. If this
// isn't set we try to use the function pointer to look up the function
// package and name.
func WithName(n string) Option {
	return func(c *config) { c.Name = n }
}

// WithValues sets some argument values that are always available to the
// function. These will also be available to any mappers when a Chain
// function is called on this func. This is useful to add some extra values
// that may be necessary to get to the types required by this function.
func WithValues(vs ...interface{}) Option {
	return func(c *config) { c.Values = vs }
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
	// Add any extra values
	values = append(values, f.Values...)

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
			// Check if the type if there are alternate values we should
			// be looking for to satisfy this. This is only in cases where
			// the arg expects multiple types.
			types := arg.Missing(values...)
			if types == nil {
				types = []Type{arg}
			}

			for _, typ := range types {
				missing[typ] = idx
			}
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
