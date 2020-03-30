// Package convert contains helpers for implementations of history.Client
// to manage the Type field on the Lookup config and convert a result to that
// type.
package convert

import (
	"reflect"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
)

// Component converts a raw input value to a requested Lookup.Type
// value using mappers, and then converts that to a slice of the resulting
// component type.
func Component(set mapper.Set, input, lookup, result interface{}) (interface{}, error) {
	// lookup will be in the form of (*Type)(nil) but we really want to
	// convert to a slice for components of (*[]Type)(nil) so this converts
	// to that.
	lookupType := reflect.TypeOf(lookup).Elem()
	lookupSlice := reflect.New(reflect.SliceOf(lookupType)).Interface()

	// Convert from the input type to our lookup type using mappers
	raw, err := set.ConvertType(input, lookupSlice)
	if err != nil {
		return nil, err
	}
	rawVal := reflect.ValueOf(raw)

	// Convert to our result type
	sliceType := reflect.TypeOf(result).Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(sliceType), rawVal.Len(), rawVal.Len())
	for i := 0; i < rawVal.Len(); i++ {
		slice.Index(i).Set(rawVal.Index(i))
	}

	return slice.Interface(), nil
}
