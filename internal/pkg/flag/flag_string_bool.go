package flag

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/posener/complete"
)

// -- StringBoolVar and stringBoolValue
type StringBoolVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    string
	Hidden     bool
	EnvVar     string
	Target     *string
	Completion complete.Predictor
	SetHook    func(val string)
}

func (f *Set) StringBoolVar(i *StringBoolVar) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		initial = v
	}

	def := ""
	if i.Default != "" {
		def = i.Default
	}

	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		Default:    def,
		EnvVar:     i.EnvVar,
		Value:      newStringBoolValue(i, initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type stringBoolValue struct {
	v      *StringBoolVar
	hidden bool
	target *string
}

func newStringBoolValue(v *StringBoolVar, def string, target *string, hidden bool) *stringBoolValue {
	*target = def
	return &stringBoolValue{
		v:      v,
		hidden: hidden,
		target: target,
	}
}

func (s *stringBoolValue) Set(val string) error {
	_, err := strconv.ParseBool(val)
	if err != nil {
		// omit the parsing error as it's only ever "invalid syntax" in favor of
		// giving users a hint. This package will expand the error to include
		// which flag had the invalid valid.
		return errors.New("please use 'true' or 'false'")
	}

	// Stay consistent with Go's acceptable values for parsing booleans, so we
	// only lowercase after the value is determined to be valid
	*s.target = strings.ToLower(val)

	if s.v.SetHook != nil {
		s.v.SetHook(val)
	}

	return nil
}

func (s *stringBoolValue) Get() interface{} { return *s.target }
func (s *stringBoolValue) String() string   { return *s.target }
func (s *stringBoolValue) Example() string  { return "true" }
func (s *stringBoolValue) Hidden() bool     { return s.hidden }
func (s *stringBoolValue) IsBoolFlag() bool { return true }
