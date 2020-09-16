package docs

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Details struct {
	Description string
	Example     string

	Input  string
	Output string

	Mappers []Mapper
}

type Mapper struct {
	Input       string
	Output      string
	Description string
}

type FieldDocs struct {
	Field    string
	Type     string
	Synopsis string
	Summary  string

	Optional bool
	Default  string
	EnvVar   string
}

type Documentation struct {
	description string
	example     string
	input       string
	output      string
	fields      map[string]*FieldDocs
	mappers     []Mapper
}

type Option func(*Documentation) error

func New(opts ...Option) (*Documentation, error) {
	var d Documentation

	d.fields = make(map[string]*FieldDocs)

	for _, opt := range opts {
		err := opt(&d)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func FromConfig(v interface{}) Option {
	return func(d *Documentation) error {
		rv := reflect.ValueOf(v).Elem()
		if rv.Kind() != reflect.Struct {
			return fmt.Errorf("invalid config type, must be struct")
		}

		t := rv.Type()

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			name, ok := f.Tag.Lookup("hcl")
			if !ok {
				return fmt.Errorf("missing hcl tag on field: %s", f.Name)
			}

			parts := strings.Split(name, ",")

			if parts[0] == "" {
				continue
			}

			field := &FieldDocs{
				Field: parts[0],
				Type:  f.Type.String(),
			}

			for _, p := range parts[1:] {
				if p == "optional" {
					field.Optional = true
				}
			}

			d.fields[parts[0]] = field
		}

		return nil
	}
}

func formatHelp(lines ...string) string {
	var sb strings.Builder

	for i, line := range lines {
		if i > 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(strings.TrimSpace(line))
	}

	return sb.String()
}

type (
	SummaryString string
	Default       string
	EnvVar        string
)

type DocOption interface {
	docOption() bool
}

func (o SummaryString) docOption() bool { return true }
func (o Default) docOption() bool       { return true }
func (o EnvVar) docOption() bool        { return true }

func Summary(in ...string) SummaryString {
	var sb strings.Builder

	for i, str := range in {
		if str == "" {
			sb.WriteByte('\n')
		}

		if i > 0 {
			sb.WriteByte(' ')
		}

		sb.WriteString(strings.TrimSpace(str))
	}

	return SummaryString(sb.String())

}

func (d *Documentation) Example(x string) {
	d.example = x
}

func (d *Documentation) Description(x string) {
	d.description = x
}

func (d *Documentation) Input(x string) {
	d.input = x
}

func (d *Documentation) Output(x string) {
	d.output = x
}

func (d *Documentation) AddMapper(input, output, description string) {
	d.mappers = append(d.mappers, Mapper{
		Input:       input,
		Output:      output,
		Description: description,
	})
}

func (d *Documentation) SetField(name, synposis string, opts ...DocOption) error {
	field, ok := d.fields[name]
	if !ok {
		field = &FieldDocs{
			Field:    name,
			Synopsis: synposis,
		}
		d.fields[name] = field
	} else {
		field.Synopsis = synposis
	}

	for _, o := range opts {
		switch v := o.(type) {
		case SummaryString:
			field.Summary = string(v)
		case Default:
			field.Default = string(v)
		case EnvVar:
			field.EnvVar = string(v)
		}
	}

	return nil
}

func (d *Documentation) OverrideField(f *FieldDocs) error {
	d.fields[f.Field] = f
	return nil
}

func (d *Documentation) Details() *Details {
	return &Details{
		Example:     d.example,
		Description: d.description,
		Input:       d.input,
		Output:      d.output,
		Mappers:     d.mappers,
	}
}

func (d *Documentation) Fields() []*FieldDocs {
	var keys []string

	for k := range d.fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var fields []*FieldDocs

	for _, k := range keys {
		fields = append(fields, d.fields[k])
	}

	return fields
}
