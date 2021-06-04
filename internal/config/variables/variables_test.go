package variables

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
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
			basecfg := HclBase{}

			err := hclsimple.DecodeFile(file, nil, &basecfg)
			require.NoError(err)

			schema, _ := gohcl.ImpliedBodySchema(&HclBase{})
			content, diag := basecfg.Body.Content(schema)
			require.False(diag.HasErrors())

			vars := &Variables{}
			for _, block := range content.Blocks {
				switch block.Type {
				case "variable":
					diag = vars.decodeVariableBlock(block)
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

func TestVariables_collectValues(t *testing.T) {
	cases := []struct {
		name        string
		file        string
		inputFiles  []string
		inputValues []*pb.Variable
		expected    Variables
		err         string
	}{
		{
			name:       "valid",
			file:       "valid.hcl",
			inputFiles: []string{filepath.Join("testdata", "values.hcl")},
			inputValues: []*pb.Variable{
				{
					Name:   "art",
					Value:  &pb.Variable_Str{Str: "gdbee"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: Variables{
				"art": &Variable{
					Values: []Value{
						{cty.DynamicVal, "default", hcl.Expression(nil), hcl.Range{}},
						{cty.StringVal("gdbee"), "file", hcl.Expression(nil), hcl.Range{}},
						{cty.StringVal("gdbee"), "cli", hcl.Expression(nil), hcl.Range{}},
					},
					Type: cty.String,
				},
			},
			err: "",
		},
		{
			name:       "undefined variable for pb.Variable value",
			file:       "valid.hcl",
			inputFiles: []string{filepath.Join("testdata", "values.hcl")},
			inputValues: []*pb.Variable{
				{
					Name:   "foo",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: Variables{},
			err:      "Undefined variable",
		},
		{
			name:       "invalid value type for pb.Variable",
			file:       "valid.hcl",
			inputFiles: []string{filepath.Join("testdata", "values.hcl")},
			inputValues: []*pb.Variable{
				{
					Name:   "is_good",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: Variables{},
			err:      "Invalid value for variable",
		},
		{
			name:        "no varfile",
			file:        "valid.hcl",
			inputFiles:  []string{filepath.Join("testdata", "nofile.hcl")},
			inputValues: []*pb.Variable{},
			expected:    Variables{},
			err:         "Given variables file testdata/nofile.hcl does not exist",
		},
		{
			name:        "input file not valid hcl",
			file:        "valid.hcl",
			inputFiles:  []string{filepath.Join("testdata", "nothcl")},
			inputValues: []*pb.Variable{},
			expected:    Variables{},
			err:         "Missing newline after argument",
		},
		{
			name:        "defintion in values file",
			file:        "valid.hcl",
			inputFiles:  []string{filepath.Join("testdata", "valid.hcl")},
			inputValues: []*pb.Variable{},
			expected:    Variables{},
			err:         "Variable declaration in a wpvars file",
		},
		{
			name:        "undefined var for file value",
			file:        "valid.hcl",
			inputFiles:  []string{filepath.Join("testdata", "undefined.hcl")},
			inputValues: []*pb.Variable{},
			expected:    Variables{},
			err:         "Undefined variable",
		},
		{
			name:        "invalid value type",
			file:        "valid.hcl",
			inputFiles:  []string{filepath.Join("testdata", "invalid_value.hcl")},
			inputValues: []*pb.Variable{},
			expected:    Variables{},
			err:         "Invalid value for variable",
		},
	}
	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			require := require.New(t)

			file := filepath.Join("testdata", tt.file)
			basecfg := HclBase{}

			err := hclsimple.DecodeFile(file, nil, &basecfg)
			require.NoError(err)

			var vs Variables
			diags := vs.DecodeVariableBlocks(basecfg.Body)
			require.False(diags.HasErrors())

			// collect values
			diags = vs.CollectInputValues(tt.inputFiles, tt.inputValues)
			if tt.err != "" {
				require.True(diags.HasErrors())
				require.Contains(diags.Error(), tt.err)
				return
			}

			require.False(diags.HasErrors())
			for k, v := range tt.expected {
				diff := cmp.Diff(v, vs[k], cmpOpts...)
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
		file     []string
		cliArgs  map[string]string
		expected []*pb.Variable
		err      string
	}{
		{
			"success",
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
		}, {
			"success",
			[]string{filepath.Join("testdata", "values.hcl")},
			map[string]string{"foo": "bar"},
			[]*pb.Variable{
				{
					Name:   "mug",
					Value:  &pb.Variable_Str{Str: "yeti"},
					Source: &pb.Variable_File_{},
				},
				{
					Name:   "art",
					Value:  &pb.Variable_Str{Str: "gdbee"},
					Source: &pb.Variable_File_{},
				},
				{
					Name:   "foo",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			"",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			vars, diags := SetJobInputVariables(tt.cliArgs, tt.file)
			require.False(diags.HasErrors())

			require.Equal(len(vars), len(tt.expected))
			require.Equal(vars, tt.expected)
			// TODO krantzinator: add a sort before comparing for equality
			for i, v := range vars {
				require.Equal(v, tt.expected[i])
			}
		})
	}
}

// func TestVariables_mergeValues(t *testing.T) {
// 	cases := []struct {

// 	}
// }

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
	cmpopts.IgnoreTypes(hcl.Range{}),
}
