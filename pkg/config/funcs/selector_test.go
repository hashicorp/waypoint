// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package funcs

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestSelectorMatch(t *testing.T) {
	tests := []struct {
		Map      map[string]string
		Selector string
		Want     cty.Value
		Err      bool
	}{
		{
			map[string]string{"env": "production"},
			"env == production",
			cty.BoolVal(true),
			false,
		},

		{
			map[string]string{"env": "production"},
			"env != production",
			cty.BoolVal(false),
			false,
		},

		{
			map[string]string{"waypoint/workspace": "foo"},
			"waypoint/workspace == foo",
			cty.BoolVal(true),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("selectormatch(%#v, %#v)", test.Map, test.Selector), func(t *testing.T) {
			// Build our map val
			mapValues := map[string]cty.Value{}
			for k, v := range test.Map {
				mapValues[k] = cty.StringVal(v)
			}

			got, err := SelectorMatch(cty.MapVal(mapValues), cty.StringVal(test.Selector))

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestSelectorLookup(t *testing.T) {
	tests := []struct {
		Map         map[string]string
		SelectorMap map[string]cty.Value
		Default     cty.Value
		Want        cty.Value
		Err         bool
	}{
		{
			map[string]string{"env": "production"},
			map[string]cty.Value{
				"env == production": cty.StringVal("prod"),
				"env == staging":    cty.StringVal("staging"),
			},
			cty.StringVal("unknown"),
			cty.StringVal("prod"),
			false,
		},

		{
			map[string]string{"env": "other"},
			map[string]cty.Value{
				"env == production": cty.StringVal("prod"),
				"env == staging":    cty.StringVal("staging"),
			},
			cty.StringVal("unknown"),
			cty.StringVal("unknown"),
			false,
		},

		{
			map[string]string{"env": "production"},
			map[string]cty.Value{
				"env == production": cty.StringVal("prod"),
				"env == staging":    cty.StringVal("staging"),
			},
			cty.BoolVal(false),
			cty.StringVal("prod"),
			true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("selectormatch(%#v, %#v)", test.Map, test.SelectorMap), func(t *testing.T) {
			// Build our map val
			mapValues := map[string]cty.Value{}
			for k, v := range test.Map {
				mapValues[k] = cty.StringVal(v)
			}

			got, err := SelectorLookup(
				cty.MapVal(mapValues),
				cty.MapVal(test.SelectorMap),
				test.Default,
			)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
