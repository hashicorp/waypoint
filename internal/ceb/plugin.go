package ceb

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/internal-shared/protomappers"
)

// callDynamicFunc is a helper to call the dynamic functions that plugins return.
//
// If error is nil, then the result is guaranteed to not be erroneous. Callers
// do NOT need to check result.Err().
func (ceb *CEB) callDynamicFunc(
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
