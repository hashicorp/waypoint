package funcs

import (
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
			nil,
			cty.NilVal,
			`vars map does not contain key "name"`,
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
