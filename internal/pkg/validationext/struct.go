package validationext

import (
	"errors"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// StructField returns a *validation.FieldRules (can be used as an arg to
// validation.ValidateStruct) that validates a nested struct value only if
// the struct value is non-nil.
//
// This is useful to apply validation to nested pointer structs all within
// a single call to validation.ValidateStruct. Otherwise, nil-checks with
// complex field rule slice building is necessary.
//
// The struct value (sv) parameter needs to be a _pointer to the struct
// field_, even if that's a pointer. So it should always be in the form
// `&s.Field`. Otherwise, you'll get a "field #N not in struct" internal
// error from ozzo-validation.
//
// Full example:
//
//	s := &Person{
//	  Address: &Address{Number: 0},
//	}
//
//	validation.ValidateStruct(&s,
//	  validation.Field(&s.Address, validation.Required),
//	  StructField(&s.Address, func() []*validation.FieldRules {
//	    return []*validation.FieldRules{
//	      validation.Field(&s.Address.Number, validation.Required),
//	    }
//	  }),
//	)
//
// In this example, the address number will be required, but will only
// be checked if the address is non-nil. Without StructField, you either
// have to do nil-checks outside of the validation call or you get a crash.
func StructField(sv interface{}, f func() []*validation.FieldRules) *validation.FieldRules {
	return validation.Field(sv, &structFieldRule{s: sv, rules: f})
}

// StructInterface is similar to StructField but validates interface-type
// fields that are non-nil and match the type t exactly. This is useful for
// "oneof" fields created by protobufs.
//
// If the function f is called, it is guaranteed that the field pointed to by
// sv is of type t. You can type assert it without ok-checking safely.
//
// See the docs for StructField for additional details on how this works.
func StructInterface(sv, t interface{}, f func() []*validation.FieldRules) *validation.FieldRules {
	return validation.Field(sv, &structFieldRule{s: sv, rules: f, iface: t})
}

// StructOneof is a special-case helper to validate struct values within
// a oneof field from a protobuf-generated struct. This behaves like
// StructInterface but automatically sets up validation directly into the
// nested oneof value. The returned fieldrules from f can be directly on the
// nested value which is implicitly required to be set.
//
// For oneof values that are NOT message types (structs) and are primitives
// like string, int, etc. then you should use StructInterface directly.
//
// Example:
//
// Given the protobuf of:
//
//	message Employee {
//	  oneof role {
//	    Eng eng = 1;
//	    Sales sales = 2;
//	  }
//	}
//
// The generated types look something lke this:
//
//	type Employee { Role Role }
//	type Role interface{}
//	type Role_Eng struct { Eng *Eng }
//	type Role_Sales struct { Sales *Sales }
//	type Eng { Language string }
//	...
//
// To validate this, you can do this:
//
//	var e *Employee
//	validation.ValidateStruct(&e,
//	  StructOneof(&e.Role, (*Eng)(nil), func() []*validation.FieldRules {
//	    v := e.Role.(*Role_Eng)
//	    return []*validation.FieldRules{
//	      validation.Field(&v.Eng.Language, validation.Required),
//	    }
//	  }),
//	)
//
// Notice how the callback sets validation on the nested `e.Role.Eng` directly.
// This helper saves a few lines of boilerplate and complicated pointer
// addressing to make this possible.
//
// The existence and non-emptiness of `e.Role.Eng` is validated automatically
// and does not need to be verified in the consumer code. This is done because
// the protobuf compiler adds these intermediate structs for type-safety
// reasons, but there is no use-case for the specific value to be nil if the
// intermediate wrapper struct is set. This results in the following cases:
//
// Valid:
//
//	Employee{ Role: nil }
//	Employee{ Role: &Role_Eng{ Eng: &Eng {} } }
//
// Invalid (a validation error is produced):
//
//	Employee{ Role: &Role_Eng{ Eng: nil } }
func StructOneof(sv, t interface{}, f func() []*validation.FieldRules) *validation.FieldRules {
	return validation.Field(sv, &structFieldRule{s: sv, rules: f, iface: t, oneof: true})
}

// StructJSONPB validates a jsonpb-encoded field (type []byte) within a struct.
// This will decode the value sv into the proto struct v and validate it with
// the rules returned by f.
//
// A validation error will be returned if jsonpb-decoding fails or if the value
// is not a byte slice.
//
// A side effect of this validation is that the field is decoded into v. After
// validation, you can continue to use this decoded value. Prior to decoding,
// we call v.Reset() so that the values are fully reset.
//
// Example:
//
//	req := struct{
//	  Employee []byte
//	}
//
//	var e pb.Employee
//	validation.ValidateStruct(&req,
//	  StructJSONPB(&req.Employee, &e, func() []*validation.FieldRules {
//	    return []*validation.FieldRules{
//	      validation.Field(&e.Name, validation.Required),
//	    }
//	  }),
//	)
func StructJSONPB(sv interface{}, v proto.Message, f func() []*validation.FieldRules) *validation.FieldRules {
	return validation.Field(sv, &structJSONPBRule{s: sv, v: v, rules: f})
}

// structJSONPBRule implements validation.Rule for StructJSONPB.
type structJSONPBRule struct {
	s     interface{}   // field pointer
	v     proto.Message // value to decode into
	rules func() []*validation.FieldRules
}

func (s *structJSONPBRule) Validate(v interface{}) error {
	// Should be bytes
	bs, ok := v.([]byte)
	if !ok {
		return errors.New("should be byte slice")
	}

	if err := protojson.Unmarshal(bs, s.v); err != nil {
		return err
	}

	// Call our rules
	return validation.ValidateStruct(
		s.v,
		s.rules()...,
	)
}

// structFieldRule implements validation.Rule. See StructField.
type structFieldRule struct {
	s     interface{}
	rules func() []*validation.FieldRules
	iface interface{}
	oneof bool
}

func (s *structFieldRule) Validate(interface{}) error {
	// The struct given to the StructField must be a pointer to that
	// struct, but we need a pointer to the field in the struct. This
	// bit of reflection checks if we have a pointer to a pointer and
	// if so, unwraps it. Otherwise we leave it as-is. This handles both
	// possible cases.
	sv := reflect.ValueOf(s.s)
	if sv.Kind() == reflect.Ptr {
		nested := reflect.Indirect(sv)
		if nested.Kind() == reflect.Ptr {
			sv = nested
		}

		// If iface is set then we want to actually validate the specific
		// type pointed to. ValidateStruct doesn't work with interface types
		// so we need to dereference it.
		if s.iface != nil {
			// This checks that the interface value is non-nil and its type
			// matches the specified concrete type directly.
			elem := nested.Elem()
			if !elem.IsValid() || elem.Type() != reflect.TypeOf(s.iface) {
				return nil
			}

			sv = elem
		}
	}
	direct := reflect.Indirect(sv)
	if !direct.IsValid() {
		return nil
	}

	// If this is a protobuf oneof, then we dereference the first field
	// in the struct
	if s.oneof {
		// We need to pointer to the first field in the struct.
		field := direct.Field(0).Addr().Interface()

		// This first field must be set. This is a requirement for all oneofs
		// and aligns with the logic of the protobuf compiler's code generation.
		if err := validation.Validate(field, validation.Required); err != nil {
			return err
		}

		// We next call ValidateStruct with our current field followed by
		// a StructField on our nested field (since protobufs inserts
		// exactly one field in a oneof value). This lets us build rules
		// on the nested value.
		return validation.ValidateStruct(
			sv.Interface(),
			StructField(
				field,
				s.rules,
			))
	}

	return validation.ValidateStruct(sv.Interface(), s.rules()...)
}
