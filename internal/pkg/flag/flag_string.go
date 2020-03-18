package flag

import (
	"os"

	"github.com/posener/complete"
)

// -- StringVar and stringValue
type StringVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    string
	Hidden     bool
	EnvVar     string
	Target     *string
	Completion complete.Predictor
}

func (f *Set) StringVar(i *StringVar) {
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
		Value:      newStringValue(initial, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type stringValue struct {
	hidden bool
	target *string
}

func newStringValue(def string, target *string, hidden bool) *stringValue {
	*target = def
	return &stringValue{
		hidden: hidden,
		target: target,
	}
}

func (s *stringValue) Set(val string) error {
	*s.target = val
	return nil
}

func (s *stringValue) Get() interface{} { return *s.target }
func (s *stringValue) String() string   { return *s.target }
func (s *stringValue) Example() string  { return "string" }
func (s *stringValue) Hidden() bool     { return s.hidden }
