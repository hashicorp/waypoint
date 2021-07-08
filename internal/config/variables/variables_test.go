package variables

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

func TestVariables_DecodeVariableBlock(t *testing.T) {
	cases := []struct {
		file string
		err  string
	}{
		{
			"valid.hcl",
			"",
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
		t.Run(tt.file, func(t *testing.T) {
			require := require.New(t)

			file := filepath.Join("testdata", tt.file)
			base := testConfig{}

			err := hclsimple.DecodeFile(file, nil, &base)
			require.NoError(err)

			schema, _ := gohcl.ImpliedBodySchema(&testConfig{})
			content, diags := base.Body.Content(schema)
			require.False(diags.HasErrors())

			vs := map[string]*Variable{}
			for _, block := range content.Blocks {
				switch block.Type {
				case "variable":
					v, decodeDiag := decodeVariableBlock(block)
					vs[block.Labels[0]] = v
					if decodeDiag.HasErrors() {
						diags = append(diags, decodeDiag...)
					}
				}
			}

			if tt.err == "" {
				require.False(diags.HasErrors())
				return
			}

			require.True(diags.HasErrors())
			require.Contains(diags.Error(), tt.err)
		})
	}
}

func TestVariables_readFileValues(t *testing.T) {
	cases := []struct {
		file     string
		expected []*pb.Variable
		err      string
	}{
		{
			file: "values.wpvars",
			expected: []*pb.Variable{
				{
					Name:   "art",
					Value:  &pb.Variable_Str{Str: "gdbee"},
					Source: &pb.Variable_Vcs{},
				},
				{
					Name:   "mug",
					Value:  &pb.Variable_Str{Str: "yeti"},
					Source: &pb.Variable_Vcs{},
				},
			},
			err: "",
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
			fv, diags := parseFileValues(fp, "vcs")

			if tt.err != "" {
				require.Contains(diags.Error(), tt.err)
				return
			}

			require.False(diags.HasErrors())
			require.Equal(len(fv), len(tt.expected))
			for _, v := range tt.expected {
				require.Contains(fv, v)
			}
		})
	}
}

func TestVariables_LoadVCSFile(t *testing.T) {
	cases := []struct {
		name     string
		expected []*pb.Variable
		err      string
	}{
		{
			name: "loads auto file only",
			expected: []*pb.Variable{
				{
					Name:   "mug",
					Value:  &pb.Variable_Str{Str: "ceramic"},
					Source: &pb.Variable_Vcs{},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			vars, diags := LoadAutoFiles("testdata")

			if tt.err != "" {
				require.True(diags.HasErrors())
				require.Contains(diags.Error(), tt.err)
				return
			}

			require.False(diags.HasErrors())
			require.ElementsMatch(vars, tt.expected)
		})
	}
}

func TestVariables_EvalInputValues(t *testing.T) {
	cases := []struct {
		name        string
		file        string
		inputValues []*pb.Variable
		expected    Values
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
			expected: Values{
				"art": &Value{
					cty.StringVal("gdbee"), "cli", hcl.Expression(nil), hcl.Range{},
				},
				"is_good": &Value{
					cty.BoolVal(false), "default", hcl.Expression(nil), hcl.Range{},
				},
				"whatdoesittaketobenumber": &Value{
					cty.NumberIntVal(1), "default", hcl.Expression(nil), hcl.Range{},
				},
			},
			err: "",
		},
		{
			name:        "complex types from default",
			file:        "list.hcl",
			inputValues: []*pb.Variable{},
			expected: Values{
				"testdata": &Value{
					stringListVal("pancakes"), "default", hcl.Expression(nil), hcl.Range{},
				},
			},
			err: "",
		},
		{
			name: "complex types from server",
			file: "list.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "testdata",
					Value:  &pb.Variable_Hcl{Hcl: "[\"waffles\"]"},
					Source: &pb.Variable_Server{},
				},
			},
			expected: Values{
				"testdata": &Value{
					stringListVal("waffles"), "server", hcl.Expression(nil), hcl.Range{},
				},
			},
			err: "",
		},
		{
			name: "complex types from cli",
			file: "list.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "testdata",
					Value:  &pb.Variable_Str{Str: "[\"waffles\"]"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: Values{
				"testdata": &Value{
					stringListVal("waffles"), "cli", hcl.Expression(nil), hcl.Range{},
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
			expected: Values{},
			err:      "Undefined variable",
		},
		{
			name: "invalid value type",
			file: "valid.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "is_good",
					Value:  &pb.Variable_Str{Str: "bar"},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: Values{},
			err:      "Invalid value for variable",
		},
		{
			name: "undefined var for file value",
			file: "undefined.hcl",
			inputValues: []*pb.Variable{
				{
					Name:   "is_good",
					Value:  &pb.Variable_Bool{Bool: true},
					Source: &pb.Variable_Cli{},
				},
			},
			expected: Values{},
			err:      "Undefined variable",
		},
		{
			name:        "no assigned or default value",
			file:        "no_default.hcl",
			inputValues: []*pb.Variable{},
			expected:    Values{},
			err:         "Unset variable",
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

			vs := map[string]*Variable{}
			for _, block := range content.Blocks {
				switch block.Type {
				case "variable":
					v, decodeDiag := decodeVariableBlock(block)
					vs[block.Labels[0]] = v
					if decodeDiag.HasErrors() {
						diags = append(diags, decodeDiag...)
					}
				}
			}
			require.False(diags.HasErrors())

			ivs, diags := EvaluateVariables(tt.inputValues, vs, hclog.New(&hclog.LoggerOptions{}))
			if tt.err != "" {
				require.True(diags.HasErrors())
				require.Contains(diags.Error(), tt.err)
				return
			}

			require.False(diags.HasErrors())
			for k, v := range tt.expected {
				diff := cmp.Diff(v, ivs[k], cmpOpts...)
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
		{
			"files",
			[]string{filepath.Join("testdata", "values.wpvars"), filepath.Join("testdata", "more_values.wpvars")},
			nil,
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
					Name:   "is_good",
					Value:  &pb.Variable_Bool{Bool: true},
					Source: &pb.Variable_File_{},
				},
				{
					Name:   "whatdoesittaketobenumber",
					Value:  &pb.Variable_Num{Num: 1},
					Source: &pb.Variable_File_{},
				},
			},
			"",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			vars, diags := LoadVariableValues(tt.cliArgs, tt.files)
			require.False(diags.HasErrors())

			require.Equal(len(vars), len(tt.expected))
			for _, v := range tt.expected {
				require.Contains(vars, v)
			}
		})
	}
}

// helper functions
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

func stringListVal(strings ...string) cty.Value {
	values := []cty.Value{}
	for _, str := range strings {
		values = append(values, cty.StringVal(str))
	}
	list, err := convert.Convert(cty.ListVal(values), cty.List(cty.String))
	if err != nil {
		panic(err)
	}
	return list
}
