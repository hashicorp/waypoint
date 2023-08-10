// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package flag

import (
	"strconv"

	"github.com/posener/complete"
)

// -- BoolPtrVar  and boolPtr
// BoolPtrVar is used to discern between values provided via CLI flags for
// "true", "false", and "not specified". Because of this there is no default
// value offered. If you want to use a boolean flag with a default, see
// flag_bool.go
type BoolPtrVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Hidden     bool
	EnvVar     string
	Target     **bool
	Completion complete.Predictor
	SetHook    func(val bool)
}

func (f *Set) BoolPtrVar(i *BoolPtrVar) {
	f.VarFlag(&VarFlag{
		Name:       i.Name,
		Aliases:    i.Aliases,
		Usage:      i.Usage,
		EnvVar:     i.EnvVar,
		Value:      newBoolPtr(i, i.Target, i.Hidden),
		Completion: i.Completion,
	})
}

type boolPtrValue struct {
	v      *BoolPtrVar
	hidden bool
	target **bool
}

func newBoolPtr(v *BoolPtrVar, target **bool, hidden bool) *boolPtrValue {
	return &boolPtrValue{
		v:      v,
		hidden: hidden,
		target: target,
	}
}

func (b *boolPtrValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*b.target = &v

	if b.v.SetHook != nil {
		b.v.SetHook(v)
	}

	return nil
}

func (b *boolPtrValue) Get() interface{} {
	if b.target != nil {
		return *b.target
	}
	return nil
}

func (b *boolPtrValue) String() string {
	if *b.target != nil {
		return strconv.FormatBool(**b.target)
	}
	return ""
}
func (b *boolPtrValue) Example() string  { return "" }
func (b *boolPtrValue) Hidden() bool     { return b.hidden }
func (b *boolPtrValue) IsBoolFlag() bool { return true }
