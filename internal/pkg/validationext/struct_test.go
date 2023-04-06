// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validationext

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/require"
)

func TestStructField(t *testing.T) {
	var called bool
	type Person struct {
		Name string
	}

	s1 := &struct {
		P Person
	}{
		P: Person{Name: "alice"},
	}

	s2 := &struct {
		P *Person
	}{
		P: &Person{Name: "alice"},
	}

	s3 := &struct {
		P *Person
	}{
		P: nil,
	}

	cases := []struct {
		Name        string
		Value       interface{}
		StructField *validation.FieldRules
		Called      bool
	}{
		{
			"non-pointer nested struct",
			s1,
			StructField(&s1.P, func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{
					validation.Field(&s1.P.Name, validation.Required),
				}
			}),
			true,
		},

		{
			"pointer nested struct, non-nil",
			s2,
			StructField(&s2.P, func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{
					validation.Field(&s2.P.Name, validation.Required),
				}
			}),
			true,
		},

		{
			"pointer nested struct, nil",
			s3,
			StructField(&s3.P, func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{
					validation.Field(&s3.P.Name, validation.Required),
				}
			}),
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Reset
			called = false

			require.NoError(validation.ValidateStruct(
				tt.Value,
				tt.StructField,
			))
			require.Equal(tt.Called, called)
		})
	}
}

func TestStructInterface(t *testing.T) {
	var called bool

	type Oneof interface{}
	type OneofA struct{ A string }
	type OneofB struct{ B int }

	s1 := &struct {
		V Oneof
	}{
		V: &OneofA{A: "hello"},
	}

	s2 := &struct {
		V Oneof
	}{
		// Unset
	}

	s3 := &struct {
		V Oneof
	}{
		V: &OneofB{B: 42},
	}

	cases := []struct {
		Name   string
		Value  interface{}
		Rules  *validation.FieldRules
		Called bool
	}{
		{
			"oneof field set",
			s1,
			StructInterface(&s1.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				v := s1.V.(*OneofA)
				return []*validation.FieldRules{
					validation.Field(&v.A, validation.Required),
				}
			}),
			true,
		},

		{
			"oneof field nil",
			s2,
			StructInterface(&s2.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{}
			}),
			false,
		},

		{
			"oneof field other type",
			s3,
			StructInterface(&s3.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{}
			}),
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Reset
			called = false

			require.NoError(validation.ValidateStruct(
				tt.Value,
				tt.Rules,
			))
			require.Equal(tt.Called, called)
		})
	}
}

func TestStructOneof(t *testing.T) {
	var called bool

	// These types mimic the types setup by a oneof in protobuf
	type A struct{ A string }
	type B struct{ B int }
	type Oneof interface{}
	type OneofA struct{ A *A }
	type OneofB struct{ B *B }

	s1 := &struct {
		V Oneof
	}{
		V: &OneofA{A: &A{A: "hello"}},
	}

	s2 := &struct {
		V Oneof
	}{
		// Unset
	}

	s3 := &struct {
		V Oneof
	}{
		V: &OneofA{A: nil},
	}

	s4 := &struct {
		V Oneof
	}{
		V: &OneofB{B: &B{B: 42}},
	}

	cases := []struct {
		Name   string
		Value  interface{}
		Rules  *validation.FieldRules
		Called bool
		Valid  bool
	}{
		{
			"oneof field set",
			s1,
			StructOneof(&s1.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				v := s1.V.(*OneofA)
				return []*validation.FieldRules{
					validation.Field(&v.A.A, validation.Required),
				}
			}),
			true,
			true,
		},

		{
			"oneof field not set",
			s2,
			StructOneof(&s2.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				v := s2.V.(*OneofA)
				return []*validation.FieldRules{
					validation.Field(&v.A.A, validation.Required),
				}
			}),
			false,
			true,
		},

		{
			"oneof field nested value nil",
			s3,
			StructOneof(&s3.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				v := s3.V.(*OneofA)
				return []*validation.FieldRules{
					validation.Field(&v.A.A, validation.Required),
				}
			}),
			false,
			false,
		},

		{
			"oneof field set to different type",
			s4,
			StructOneof(&s4.V, (*OneofA)(nil), func() []*validation.FieldRules {
				called = true
				v := s4.V.(*OneofA)
				return []*validation.FieldRules{
					validation.Field(&v.A.A, validation.Required),
				}
			}),
			false, // should not be called!
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Reset
			called = false

			err := validation.ValidateStruct(
				tt.Value,
				tt.Rules,
			)
			require.Equal(tt.Called, called)
			if tt.Valid {
				require.NoError(err)
			} else {
				require.Error(err)
			}
		})
	}
}

/*
func TestStructJSONPB(t *testing.T) {
	m := func(p *PersonPB) []byte {
		var m jsonpb.Marshaler
		var buf bytes.Buffer
		require.NoError(t, m.Marshal(&buf, p))
		return buf.Bytes()
	}

	s1 := &struct {
		P []byte
	}{
		P: m(&PersonPB{Name: "bob", Pronoun: "they"}),
	}
	s2 := &struct {
		P []byte
	}{
		P: m(&PersonPB{Name: "", Pronoun: "they"}),
	}

	var called bool
	var p PersonPB
	cases := []struct {
		Name        string
		Value       interface{}
		StructField *validation.FieldRules
		Called      bool
		Error       string
	}{
		{
			"valid",
			s1,
			StructJSONPB(&s1.P, &p, func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{
					validation.Field(&p.Name, validation.Required),
				}
			}),
			true,
			"",
		},

		{
			"error",
			s2,
			StructJSONPB(&s2.P, &p, func() []*validation.FieldRules {
				called = true
				return []*validation.FieldRules{
					validation.Field(&p.Name, validation.Required),
				}
			}),
			true,
			"name: cannot be blank",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Reset
			called = false
			p = PersonPB{}

			err := validation.ValidateStruct(
				tt.Value,
				tt.StructField,
			)
			require.Equal(tt.Called, called)
			if tt.Error == "" {
				require.NoError(err)
				return
			}
			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
*/

// PersonPB is a proto.Message used for testing.
type PersonPB struct {
	Name    string `protobuf:"bytes,2,opt,name=name,json=name,proto3" json:"name,omitempty"`
	Pronoun string `protobuf:"bytes,3,opt,name=pronoun,json=pronoun,proto3" json:"pronoun,omitempty"`
}

// Implement proto.Message with dummy implementation
func (p *PersonPB) ProtoMessage()  {}
func (p *PersonPB) Reset()         { *p = PersonPB{} }
func (p *PersonPB) String() string { return "" }
