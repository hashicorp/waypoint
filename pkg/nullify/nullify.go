// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package nullify takes any structure as input and nullifies any matching
// pointer types as registered on the nullifier.
package nullify

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/reflectwalk"
)

// Nullify takes an input and nullifies all the types t. The types t should
// be pointer types, i.e. (*Foo)(nil).
func Nullify(input interface{}, ts ...interface{}) error {
	types := map[reflect.Type]reflect.Value{}
	for _, raw := range ts {
		t := reflect.TypeOf(raw)
		if t.Kind() != reflect.Ptr {
			return fmt.Errorf("type must be a pointer, got: %s", t)
		}

		types[t.Elem()] = reflect.Zero(t)
	}

	return reflectwalk.Walk(input, &nullifier{
		Types: types,
	})
}

type nullifier struct {
	Types map[reflect.Type]reflect.Value
}

func (n *nullifier) Struct(v reflect.Value) error {
	return nil
}

func (n *nullifier) StructField(sf reflect.StructField, v reflect.Value) error {
	if sf.PkgPath != "" {
		return reflectwalk.SkipEntry
	}

	// If is null already or not addressable, can't do anything.
	// NOTE the "not address" or CanSet check is a result of many reasons why
	// this might not work. The user must pass in an addressable value, meaning
	// a pointer, non-const, etc.
	if !v.IsValid() {
		return nil
	}
	if !v.CanSet() {
		return nil
	}

	// We want to find a pointer. If it is already nil then we're good.
	vt := v.Type()
	if vt.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}

	// If this isn't a type we care about, ignore.
	zero, ok := n.Types[vt.Elem()]
	if !ok {
		return nil
	}

	// Set it to nil.
	v.Set(zero)
	return nil
}

var (
	_ reflectwalk.StructWalker = (*nullifier)(nil)
)
