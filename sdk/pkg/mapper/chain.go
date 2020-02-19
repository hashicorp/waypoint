package mapper

import (
	"fmt"
	"reflect"
	"strings"
)

// NOTE(mitchellh): The whole algorithm below is sub-optimal in many ways:
// we use too much state, we duplicate processing work, etc. We can improve
// this as needed since the tests are written and are high level.

// Chain is similar to Prepare or Call but takes a list of Funcs that
// can be called as intermediaries to convert parameters to the expected
// parameters for this func. The result is a "Chain" of function calls.
func (f *Func) Chain(mappers []*Func, values ...interface{}) (*Chain, error) {
	// First, we need to determine what we're missing for our func.
	missing := make(map[Type]int)
	f.args(values, missing)

	// If we're not missing anything then short-circuit the whole thing
	if len(missing) == 0 {
		return &Chain{funcs: []*Func{f}, values: values}, nil
	}
	missing = nil // We don't need this anymore

	// Build a map of what our functions all provide
	mapperByOut := make(map[interface{}][]*Func)
	for _, m := range mappers {
		mapperByOut[m.Out.Key()] = append(mapperByOut[m.Out.Key()], m)
	}

	// Build our chain
	chain, err := f.chain(
		values,
		mapperByOut,
		make([]*Func, 0, 1),
		make(map[*Func]struct{}),
		make(map[*Func]struct{}),
	)
	if err != nil {
		return nil, err
	}

	return &Chain{funcs: chain, values: values}, nil
}

// chain is the internal recursive functions called on functions to build
// up the chain.
func (f *Func) chain(
	values []interface{},
	mapperByOut map[interface{}][]*Func, // mappers by output type
	chain []*Func, // chain so far
	chainSet map[*Func]struct{}, // set of functions we're calling so far
	pendingSet map[*Func]struct{}, // stack of functions that aren't yet satisfied
) ([]*Func, error) {
	missing := make(map[Type]int)
	f.args(values, missing)

	// If we have no missing values, we're satisfied
	if len(missing) == 0 {
		chainSet[f] = struct{}{}
		return append(chain, f), nil
	}

	// Add ourselves immediately to the pending set since we're no longer valid
	pendingSet[f] = struct{}{}
	defer delete(pendingSet, f)

MISSING_LOOP:
	// Go through each missing value and look for a func that will produce it
	for t, _ := range missing {
		ms := mapperByOut[t.Key()]
		if len(ms) > 0 {
			// See if we call any of these mappers already. If we do, then
			// we're satisfied by that and we can skip this missing value.
			for _, m := range ms {
				if _, ok := chainSet[m]; ok {
					continue MISSING_LOOP
				}
			}

			// Not satisfied yet so we go through each mapper and try to find
			// one that can be satisfied by our inputs.
			for _, m := range ms {
				// Skip any mappers in the pending set, since those are still
				// trying to be satisfied and if we tried to call it we'd
				// loop.
				if _, ok := pendingSet[m]; ok {
					continue
				}

				nextChain, err := m.chain(values, mapperByOut, chain, chainSet, pendingSet)
				if err == nil {
					// Satisfied!
					chain = nextChain
					continue MISSING_LOOP
				}
			}
		}

		return nil, fmt.Errorf("unable to map to %s", t.String())
	}

	return append(chain, f), nil
}

// Chain represents a chain of functions that need to be called to build
// values to satisfy the inputs of the subsequent functions.
type Chain struct {
	// funcs is an ordered list of functions that need to be called.
	funcs []*Func

	// values is the list of values we have
	values []interface{}
}

// Call calls all the functions in the chain and returns the first error
// or final result.
func (c *Chain) Call() (interface{}, error) {
	var result interface{}
	var err error
	for _, f := range c.funcs {
		result, err = f.Prepare(c.values...).Call()
		if err != nil {
			return nil, err
		}

		v := reflect.ValueOf(result)
		if v.IsValid() {
			c.values = append(c.values, result)
		}
	}

	return result, nil
}

// String implements Stringer and outputs a human-friendly description
// of the call chain that this represents.
func (c *Chain) String() string {
	ss := make([]string, len(c.funcs))
	for i, f := range c.funcs {
		ss[i] = f.String()
	}

	return strings.Join(ss, " => ")
}
