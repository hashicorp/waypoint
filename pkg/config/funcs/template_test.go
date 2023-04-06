// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestTemplateString(t *testing.T) {
	tests := []struct {
		Name     string
		Template cty.Value
		Vars     []cty.Value
		Want     cty.Value
		Err      string
	}{
		{
			"string",
			cty.StringVal("Hello World"),
			nil,
			cty.StringVal("Hello World"),
			``,
		},
		{
			"template",
			cty.StringVal("Hello, ${name}!"),
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"name": cty.StringVal("Jodie"),
				}),
			},
			cty.StringVal("Hello, Jodie!"),
			``,
		},
		{
			"template with object",
			cty.StringVal("Hello, ${name}!"),
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("Jimbo"),
				}),
			},
			cty.StringVal("Hello, Jimbo!"),
			``,
		},
		{
			"missing variable",
			cty.StringVal("Hello, ${name}!"),
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"not_name": cty.StringVal("Jodie"),
				})},
			cty.NilVal,
			`vars map does not contain key "name"`,
		},
		{
			"no variables at all",
			cty.StringVal("Hello, ${name}!"),
			nil, // must blame the template for the missing variable in this case
			cty.NilVal,
			`but this call has no vars map`,
		},
		{
			"parent value",
			cty.StringVal("Hello, ${animal}!"),
			nil,
			cty.StringVal(`Hello, dog!`),
			"",
		},
		{
			"recursive",
			cty.StringVal(`Hello ${templatestring("foo")}!`),
			nil,
			cty.NilVal,
			`cannot call templatestring from inside`,
		},
	}

	parentCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"animal": cty.StringVal("dog"),
		},
	}
	templateFns := MakeTemplateFuncs(parentCtx)
	templateFn := templateFns["templatestring"]

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			got, err := templateFn.Call(append([]cty.Value{tt.Template}, tt.Vars...))
			if tt.Err != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
				return
			}
			require.NoError(err)
			require.Equal(tt.Want, got)
		})
	}
}

func TestTemplateFile(t *testing.T) {
	tests := []struct {
		Path cty.Value
		Vars []cty.Value
		Want cty.Value
		Err  string
	}{
		{
			cty.StringVal("testdata/filesystem/hello.txt"),
			nil,
			cty.StringVal("Hello World"),
			``,
		},
		{
			cty.StringVal("testdata/filesystem/icon.png"),
			nil,
			cty.NilVal,
			`are not valid UTF-8; use the filebase64 function to obtain the Base64 encoded contents or the other file functions (e.g. filemd5, filesha256) to obtain file hashing results instead`,
		},
		{
			cty.StringVal("testdata/filesystem/missing"),
			nil,
			cty.NilVal,
			`no file exists`,
		},
		{
			cty.StringVal("testdata/filesystem/hello.tmpl"),
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"name": cty.StringVal("Jodie"),
				}),
			},
			cty.StringVal("Hello, Jodie!"),
			``,
		},
		{
			cty.StringVal("testdata/filesystem/hello.tmpl"),
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"name!": cty.StringVal("Jodie"),
				}),
			},
			cty.NilVal,
			`invalid template variable name "name!": must start with a letter, followed by zero or more letters, digits, and underscores`,
		},
		{
			cty.StringVal("testdata/filesystem/hello.tmpl"),
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("Jimbo"),
				}),
			},
			cty.StringVal("Hello, Jimbo!"),
			``,
		},
		{
			cty.StringVal("testdata/filesystem/hello.tmpl"),
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"not_name": cty.StringVal("Jodie"),
				}),
			},
			cty.NilVal,
			`vars map does not contain key "name"`,
		},
		{
			cty.StringVal("testdata/filesystem/hello.tmpl"),
			nil, // must blame the template for the missing variable in this case
			cty.NilVal,
			`but this call has no vars map`,
		},
		{
			cty.StringVal("testdata/filesystem/func.tmpl"),
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"list": cty.ListVal([]cty.Value{
						cty.StringVal("a"),
						cty.StringVal("b"),
						cty.StringVal("c"),
					}),
				}),
			},
			cty.StringVal("The items are a, b, c"),
			``,
		},
		{
			cty.StringVal("testdata/filesystem/recursive.tmpl"),
			[]cty.Value{
				cty.MapValEmpty(cty.String),
			},
			cty.NilVal,
			`cannot call templatefile from inside a template function.`,
		},
		{
			cty.StringVal("testdata/filesystem/list.tmpl"),
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"list": cty.ListVal([]cty.Value{
						cty.StringVal("a"),
						cty.StringVal("b"),
						cty.StringVal("c"),
					}),
				}),
			},
			cty.StringVal("- a\n- b\n- c\n"),
			``,
		},
		{
			cty.StringVal("testdata/filesystem/list.tmpl"),
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"list": cty.True,
				}),
			},
			cty.NilVal,
			`over non-iterable value; A value of type bool cannot be used as the collection in a 'for' expression.`,
		},
		{
			cty.StringVal("testdata/filesystem/bare.tmpl"),
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"val": cty.True,
				}),
			},
			cty.StringVal("true"),
			``,
		},
	}

	parentCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"animal": cty.StringVal("dog"),
		},
		Functions: Stdlib(),
	}
	templateFns := MakeTemplateFuncs(parentCtx)
	templateFn := templateFns["templatefile"]

	for _, tt := range tests {
		t.Run(tt.Path.AsString(), func(t *testing.T) {
			require := require.New(t)

			abs, err := filepath.Abs(tt.Path.AsString())
			require.NoError(err)
			tt.Path = cty.StringVal(abs)

			got, err := templateFn.Call(append([]cty.Value{tt.Path}, tt.Vars...))
			if tt.Err != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
				return
			}
			require.NoError(err)

			data, err := ioutil.ReadFile(got.AsString())
			require.NoError(err)

			require.Equal(tt.Want, cty.StringVal(string(data)))
		})
	}
}
