// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package funcs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestJsonnetFile(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Opts cty.Value
		Err  string
	}{
		{
			cty.StringVal("testdata/jsonnet/hello.jsonnet"),
			cty.ObjectVal(nil),
			``,
		},

		{
			cty.StringVal("testdata/jsonnet/imports.jsonnet"),
			cty.ObjectVal(nil),
			``,
		},

		{
			cty.StringVal("testdata/jsonnet/top-level-ext.jsonnet"),
			cty.ObjectVal(map[string]cty.Value{
				"ext_vars": cty.MapVal(map[string]cty.Value{
					"prefix": cty.StringVal("Happy Hour "),
				}),

				"ext_code": cty.MapVal(map[string]cty.Value{
					"brunch": cty.StringVal("true"),
				}),
			}),
			``,
		},

		{
			cty.StringVal("testdata/jsonnet/top-level-ext.jsonnet"),
			// with a map
			cty.MapVal(map[string]cty.Value{
				"ext_vars": cty.MapVal(map[string]cty.Value{
					"prefix": cty.StringVal("Happy Hour "),
				}),

				"ext_code": cty.MapVal(map[string]cty.Value{
					"brunch": cty.StringVal("true"),
				}),
			}),
			``,
		},

		{
			cty.StringVal("testdata/jsonnet/top-level-tla.jsonnet"),
			cty.ObjectVal(map[string]cty.Value{
				"tla_vars": cty.MapVal(map[string]cty.Value{
					"prefix": cty.StringVal("Happy Hour "),
				}),

				"tla_code": cty.MapVal(map[string]cty.Value{
					"brunch": cty.StringVal("true"),
				}),
			}),
			``,
		},

		// extra arguments are just ignored
		{
			cty.StringVal("testdata/jsonnet/hello.jsonnet"),
			cty.ObjectVal(map[string]cty.Value{
				"what": cty.MapVal(map[string]cty.Value{
					"foo": cty.StringVal("bar"),
				}),
			}),
			``,
		},

		// TLA vars with invalid type
		{
			cty.StringVal("testdata/jsonnet/hello.jsonnet"),
			cty.ObjectVal(map[string]cty.Value{
				"tla_vars": cty.MapVal(map[string]cty.Value{
					"foo": cty.BoolVal(true),
				}),
			}),
			`string types`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Path.AsString(), func(t *testing.T) {
			require := require.New(t)

			abs, err := filepath.Abs(tt.Path.AsString())
			require.NoError(err)
			tt.Path = cty.StringVal(abs)

			got, err := JsonnetFileFunc.Call([]cty.Value{
				tt.Path,
				tt.Opts,
			})
			if tt.Err != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
				return
			}
			require.NoError(err)

			// Ensure that our file ends in ".json"
			path := got.AsString()
			require.Equal(filepath.Ext(path), ".json")

			data, err := ioutil.ReadFile(path)
			require.NoError(err)

			const outSuffix = ".out"
			g := goldie.New(t,
				goldie.WithFixtureDir(filepath.Join("testdata", "jsonnet")),
				goldie.WithNameSuffix(outSuffix),
			)
			g.Assert(t, filepath.Base(tt.Path.AsString()), data)
		})
	}
}

func TestJsonnetDir(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Opts cty.Value
		Err  string
	}{
		{
			cty.StringVal("testdata/jsonnet/dir"),
			cty.ObjectVal(map[string]cty.Value{
				"ext_vars": cty.MapVal(map[string]cty.Value{
					"prefix": cty.StringVal("Happy Hour "),
				}),

				"ext_code": cty.MapVal(map[string]cty.Value{
					"brunch": cty.StringVal("true"),
				}),
			}),
			``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Path.AsString(), func(t *testing.T) {
			require := require.New(t)

			abs, err := filepath.Abs(tt.Path.AsString())
			require.NoError(err)
			tt.Path = cty.StringVal(abs)

			result, err := JsonnetDirFunc.Call([]cty.Value{
				tt.Path,
				tt.Opts,
			})
			if tt.Err != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
				return
			}
			require.NoError(err)

			// Go through each file and validate that we have a match.
			path := result.AsString()
			require.NoError(filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				data, err := ioutil.ReadFile(path)
				require.NoError(err)

				const outSuffix = ".out"
				g := goldie.New(t,
					goldie.WithFixtureDir(tt.Path.AsString()),
					goldie.WithNameSuffix(outSuffix),
				)
				g.Assert(t, filepath.Base(path), data)

				return nil
			}))
		})
	}
}
