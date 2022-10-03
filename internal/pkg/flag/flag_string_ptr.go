package flag

import (
	"github.com/posener/complete"
)

// -- StringPtrVar and stringValue
type StringPtrVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    string
	Hidden     bool
	EnvVar     string
	Target     **string
	Completion complete.Predictor
	SetHook    func(val string)
}

func (f *Set) StringPtrVar(i *StringPtrVar) {
	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Default:    i.Default,
		Usage:      i.Usage,
		EnvVar:     i.EnvVar,
		Value:      newStringPtr(i, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type stringPtrValue struct {
	v      *StringPtrVar
	hidden bool
	target **string
}

func newStringPtr(v *StringPtrVar, target **string, hidden bool) *stringPtrValue {
	return &stringPtrValue{
		v:      v,
		hidden: hidden,
		target: target,
	}
}

func (s *stringPtrValue) Set(val string) error {
	*s.target = &val

	if s.v.SetHook != nil {
		s.v.SetHook(val)
	}

	return nil
}

func (s *stringPtrValue) String() string {
	if *s.target != nil {
		return **s.target
	}
	return s.v.Default
}

func (s *stringPtrValue) Get() interface{} { return s.target }
func (s *stringPtrValue) Example() string  { return "string" }
func (s *stringPtrValue) Hidden() bool     { return s.hidden }
