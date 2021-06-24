package variables

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestVariables_decode(t *testing.T) {
	// TODO krantzinator: this can probably move under just validate, and
	// then Decode calls the validate function and if it passes, saves those in
	// *variables
	cases := []struct {
		File string
		Err  string
	}{
		{
			"valid.hcl",
			"",
		},
		{
			"duplicate_def.hcl",
			"Duplicate variable",
		},
		{
			"invalid_name.hcl",
			"Invalid variable name",
		},
		{
			"invalid_def.hcl",
			"Invalid default value",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			file := filepath.Join("testdata", tt.File)
			base := testConfig{}

			err := hclsimple.DecodeFile(file, nil, &base)
			require.NoError(err)

			schema, _ := gohcl.ImpliedBodySchema(&testConfig{})
			content, diag := base.Body.Content(schema)
			require.False(diag.HasErrors())

			vars := &InputVars{}
			for _, block := range content.Blocks {
				switch block.Type {
				case "variable":
					diag = vars.DecodeVariableBlock(block)
				}
			}

			if tt.Err == "" {
				require.False(diag.HasErrors())
				return
			}

			require.True(diag.HasErrors())
			require.Contains(diag.Error(), tt.Err)
		})
	}
}

func TestVariables_readFileValues(t *testing.T) {
	cases := []struct {
		file string
		err  string
	}{
		{
			file: "values.wpvars",
			err:  "",
		},
		{
			file: "nofile.wpvars",
			err:  "Given variables file testdata/nofile.wpvars does not exist",
		},
		{
			file: "nothcl",
			err:  "Missing newline after argument",
		},
		{
			file: "valid.hcl",
			err:  "Variable declaration in a wpvars file",
		},
	}
	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			require := require.New(t)

			fp := filepath.Join("testdata", tt.file)
			_, diags := readFileValues(fp)

			if tt.err != "" {
				require.Contains(diags.Error(), tt.err)
				return
			}

			require.False(diags.HasErrors())
		})
	}
}

func TestVariables_collectValues(t *testing.T) {
	cases := []struct {
		name        string
		file        string
		inputValues []*pb.Variable
		expected    InputVars
		err         string
	}{
		{
			name: "valid",
			file: "valid.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "art",
					Value:  &pb.Variable_Str{Str: "gdbee"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: InputVars{
				"art": &InputVar{
					Name: "art",
					Values: []Value{
						{cty.StringVal("gdbee"), Source{"cli", 5}, hcl.Expression(nil), hcl.Range{}},
						{cty.StringVal("gdbee"), Source{"vcs", 2}, hcl.Expression(nil), hcl.Range{}},
						{cty.NullVal(cty.String), Source{"default", 0}, hcl.Expression(nil), hcl.Range{}},
					},
					Type: cty.String,
				},
			},
			err: "",
		},
		{
			name: "undefined variable for pb.Variable value",
			file: "valid.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "foo",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: InputVars{},
			err:      "Undefined variable",
		},
		{
			name: "invalid value type for pb.Variable",
			file: "valid.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "is_good",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: InputVars{},
			err:      "Invalid value for variable",
		},
		{
			name:        "undefined var for file value",
			file:        "undefined.hcl",
			inputValues: []*pb.Variable{},
			expected:    InputVars{},
			err:         "Undefined variable",
		},
		{
			name:        "invalid value type",
			file:        "invalid_value.hcl",
			inputValues: []*pb.Variable{},
			expected:    InputVars{},
			err:         "Invalid value for variable",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			file := filepath.Join("testdata", tt.file)
			base := testConfig{}

			err := hclsimple.DecodeFile(file, nil, &base)
			require.NoError(err)

			schema, _ := gohcl.ImpliedBodySchema(&testConfig{})
			content, diags := base.Body.Content(schema)
			require.False(diags.HasErrors())

			vars := &InputVars{}
			for _, block := range content.Blocks {
				switch block.Type {
				case "variable":
					diags = vars.DecodeVariableBlock(block)
				}
			}
			require.False(diags.HasErrors())

			// collect values
			diags = vars.CollectInputValues("testdata", tt.inputValues)
			if tt.err != "" {
				require.True(diags.HasErrors())
				require.Contains(diags.Error(), tt.err)
				return
			}

			require.False(diags.HasErrors())
			for k, v := range tt.expected {
				diff := cmp.Diff(v, (*vars)[k], cmpOpts...)
				if diff != "" {
					t.Fatalf("Expected variables differed from actual: %s", diff)
				}
			}
		})
	}
}

func TestVariables_SetJobInputVariables(t *testing.T) {
	cases := []struct {
		name     string
		files    []string
		cliArgs  map[string]string
		expected []*pb.Variable
		err      string
	}{
		{
			"cli args",
			[]string{""},
			map[string]string{"foo": "bar"},
			[]*pb.Variable{
				{
					Name:   "foo",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			"",
		},
		// {
		// 	"files",
		// 	[]string{filepath.Join("testdata", "values.wpvars"), filepath.Join("testdata", "more_values.wpvars")},
		// 	nil,
		// 	[]*pb.Variable{
		// 		{
		// 			Name:   "mug",
		// 			Value:  &pb.Variable_Str{Str: "yeti"},
		// 			Source: &pb.Variable_File_{},
		// 		},
		// 		{
		// 			Name:   "art",
		// 			Value:  &pb.Variable_Str{Str: "gdbee"},
		// 			Source: &pb.Variable_File_{},
		// 		},
		// 		{
		// 			Name:   "is_good",
		// 			Value:  &pb.Variable_Str{Str: "true"},
		// 			Source: &pb.Variable_File_{},
		// 		},
		// 	},
		// 	"",
		// },
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			vars, diags := SetJobInputVariables(tt.cliArgs, tt.files)
			require.False(diags.HasErrors())

			require.Equal(len(vars), len(tt.expected))
			for _, v := range tt.expected {
				require.Contains(vars, v)
			}
		})
	}
}

var ctyValueComparer = cmp.Comparer(func(x, y cty.Value) bool {
	return x.RawEquals(y)
})

var ctyTypeComparer = cmp.Comparer(func(x, y cty.Type) bool {
	if x == cty.NilType && y == cty.NilType {
		return true
	}
	if x == cty.NilType || y == cty.NilType {
		return false
	}
	return x.Equals(y)
})

var cmpOpts = []cmp.Option{
	ctyValueComparer,
	ctyTypeComparer,
	cmpopts.IgnoreInterfaces(struct{ hcl.Expression }{}),
	cmpopts.IgnoreTypes(hclsyntax.TemplateExpr{}),
	cmpopts.IgnoreTypes(hcl.Range{}),
}

type testConfig struct {
	Variables []*HclVariable `hcl:"variable,block"`
	Body      hcl.Body       `hcl:",body"`
}
