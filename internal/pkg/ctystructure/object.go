package ctystructure

import (
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// Object takes a map[string]interface{} and converts it to a cty object value.
// The map[string]interface{} is expected to only have primitives and container
// types. Lists must be homogeneous types.
func Object(v map[string]interface{}) (cty.Value, error) {
	path := make(cty.Path, 0)
	return toValueObject(reflect.ValueOf(v), path)
}

func toValue(val reflect.Value, path cty.Path) (cty.Value, error) {
	val = unwrapPointer(val)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cty.NumberIntVal(val.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cty.NumberUIntVal(val.Uint()), nil

	case reflect.Float32, reflect.Float64:
		return cty.NumberFloatVal(val.Float()), nil

	case reflect.String:
		return cty.StringVal(val.String()), nil

	case reflect.Bool:
		return cty.BoolVal(val.Bool()), nil

	case reflect.Array, reflect.Slice:
		return toValueList(val, path)

	case reflect.Map:
		return toValueObject(val, path)

	case reflect.Invalid:
		return cty.NilVal, nil

	default:
		return cty.NilVal, path.NewErrorf(
			"unknown kind: %s", val.Kind().String())
	}
}

func toValueList(val reflect.Value, path cty.Path) (cty.Value, error) {
	// While we work on our elements we'll temporarily grow
	// path to give us a place to put our index step.
	path = append(path, cty.PathStep(nil))
	lastPath := &path[len(path)-1]

	result := make([]cty.Value, val.Len())
	for i := range result {
		var err error
		*lastPath = cty.IndexStep{Key: cty.NumberIntVal(int64(i))}

		result[i], err = toValue(val.Index(i), path)
		if err != nil {
			return cty.NilVal, err
		}

		if i > 0 && result[i].Type() != result[0].Type() {
			return cty.NilVal, path.NewErrorf(
				"all elements in a list must be the same type")
		}
	}

	return cty.ListVal(result), nil
}

func toValueObject(val reflect.Value, path cty.Path) (cty.Value, error) {
	// While we work on our elements we'll temporarily grow
	// path to give us a place to put our index step.
	path = append(path, cty.PathStep(nil))
	lastPath := &path[len(path)-1]

	attrs := map[string]cty.Value{}
	for _, kv := range val.MapKeys() {
		*lastPath = cty.IndexStep{Key: cty.StringVal(kv.String())}

		value, err := toValue(val.MapIndex(kv), path)
		if err != nil {
			return cty.NilVal, err
		}

		attrs[kv.String()] = value
	}

	return cty.ObjectVal(attrs), nil
}

// unwrapPointer is a helper for dealing with Go pointers. It has three
// possible outcomes:
//
//   - Given value isn't a pointer, so it's just returned as-is.
//   - Given value is a non-nil pointer, in which case it is dereferenced
//     and the result returned.
//   - Given value is a nil pointer, in which case an invalid value is returned.
//
// For nested pointer types, like **int, they are all dereferenced in turn
// until a non-pointer value is found, or until a nil pointer is encountered.
func unwrapPointer(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return reflect.Value{}
		}

		val = val.Elem()
	}

	return val
}
