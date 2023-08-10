// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package validationext

import (
	"errors"
	"reflect"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Note: this is copied from ozzo-validation, but modified to not unwrap
// pointers. Since our projects use protocol buffers, it is not safe to unwrap
// a pointer since that forces a copy of things such as locks.

// Each returns a validation rule that loops through an iterable (map, slice or array)
// and validates each value inside with the provided rules.
// An empty iterable is considered valid. Use the Required rule to make sure the iterable is not empty.
func Each(rules ...validation.Rule) EachRule {
	return EachRule{
		rules: rules,
	}
}

// EachRule is a validation rule that validates elements in a map/slice/array using the specified list of rules.
type EachRule struct {
	rules []validation.Rule
}

// Validate loops through the given iterable and calls the Ozzo Validate() method for each value.
func (r EachRule) Validate(value interface{}) error {
	errs := validation.Errors{}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			err := validation.Validate(v.MapIndex(k).Interface(), r.rules...)
			if err != nil {
				errs[r.getString(k)] = err
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			err := validation.Validate(v.Index(i).Interface(), r.rules...)
			if err != nil {
				errs[strconv.Itoa(i)] = err
			}
		}
	default:
		return errors.New("must be an iterable (map, slice or array)")
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (r EachRule) getString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return ""
		}
		return value.Elem().String()
	default:
		return value.String()
	}
}
