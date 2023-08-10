// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/internal-shared/protomappers"
)

// CallDynamicFunc is a helper to call the dynamic functions that Waypoint
// plugins return, i.e. a `DeployFunc`.
//
// The value f must be either a function pointer or an *argmapper.Func directly.
// The called function will always have access to the given logger and the logger
// will also be used by argmapper.
//
// If f is nil then (nil, nil) is returned. This is for prior compatibility
// reasons so please always check the return value even if err is nil.
//
// If error is nil, then the result is guaranteed to not be erroneous. Callers
// do NOT need to check result.Err().
func CallDynamicFunc(
	log hclog.Logger,
	f interface{},
	args ...argmapper.Arg,
) (*argmapper.Result, error) {
	if f == nil {
		return nil, nil
	}

	// Get our function.
	rawFunc, ok := f.(*argmapper.Func)
	if !ok {
		var err error
		rawFunc, err = argmapper.NewFunc(f, argmapper.Logger(log))
		if err != nil {
			return nil, err
		}
	}

	// Build our mappers
	mappers, err := argmapper.NewFuncList(protomappers.All, argmapper.Logger(log))
	if err != nil {
		return nil, err
	}

	// Call it
	result := rawFunc.Call(append(args,
		argmapper.ConverterFunc(mappers...),
		argmapper.Typed(log),
	)...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}
