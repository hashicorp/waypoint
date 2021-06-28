package variables

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

// A consistent detail message for all "not a valid identifier" diagnostics.
var (
	badIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."

	// Variable value sources
	// Highest precedence is the final value used
	sourceDefault = "default"
	sourceServer  = "server"
	sourceVCS     = "vcs"
	sourceEnv     = "env"
	sourceFile    = "file"
	sourceCLI     = "cli"
)

// InputValue contain the value of the variable along with associated metada,
// including the source it was set from: cli, file, env, vcs, server/ui
type InputValue struct {
	Value  cty.Value
	Source string
	Expr   hcl.Expression
	// The location of the variable value if the value was provided
	// from a file
	Range hcl.Range
}

type Variable struct {
	Name string

	Default *InputValue

	// Cty Type of the variable. If the default value or a collected value is
	// not of this type nor can be converted to this type an error diagnostic
	// will show up. This allows us to assume that values are valid.
	//
	// When a default value - and no type - is passed into the variable
	// declaration, the type of the default variable will be used.
	Type cty.Type

	// Description of the variable
	Description string

	// The location of the variable definition block in the waypoint.hcl
	Range hcl.Range
}

type HclVariable struct {
	Name        string         `hcl:",label"`
	Default     cty.Value      `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
	Description string         `hcl:"description,optional"`
}

// InputValues are used when parsing the Config, to set default values from
// the waypoint.hcl and bring in the values from the job and the server/VCS
// for eventual precedence sorting and setting on the EvalContext
type InputValues map[string]*InputValue

// TODO krantzinator - use implied body scheme instead?
var variableBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "default",
		},
		{
			Name: "type",
		},
		{
			Name: "description",
		},
	},
}

// decodeVariableBlock first validates the variable block, and then sets
// the decoded variables with their default value on *InputVars

// parseVarBlock validates each part of the variable block, building out
// a defined *Variable definition
func DecodeVariableBlock(block *hcl.Block) (*Variable, hcl.Diagnostics) {
	name := block.Labels[0]
	v := Variable{
		Name:  name,
		Range: block.DefRange,
	}

	content, diags := block.Body.Content(variableBlockSchema)

	if !hclsyntax.ValidIdentifier(name) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid variable name",
			Detail:   badIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		valDiags := gohcl.DecodeExpression(attr.Expr, nil, &v.Description)
		diags = append(diags, valDiags...)
		if diags.HasErrors() {
			return nil, diags
		}
	}

	if attr, ok := content.Attributes["type"]; ok {
		t, moreDiags := typeexpr.Type(attr.Expr)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags
		}
		v.Type = t
	}

	if attr, exists := content.Attributes["default"]; exists {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)
		if diags.HasErrors() {
			return nil, diags
		}
		// Convert the default to the expected type so we can catch invalid
		// defaults early and allow later code to assume validity.
		// Note that this depends on us having already processed any "type"
		// attribute above.
		if v.Type != cty.NilType {
			var err error
			val, err = convert.Convert(val, v.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Invalid default value for variable %q", name),
					Detail:   fmt.Sprintf("This default value is not compatible with the variable's type constraint: %s.", err),
					Subject:  attr.Expr.Range().Ptr(),
				})
				val = cty.DynamicVal
			}
		}

		v.Default = &InputValue{
			Source: sourceDefault,
			Value:  val,
		}

		// It's possible no type attribute was assigned so lets make sure we
		// have a valid type otherwise there could be issues parsing the value.
		if v.Type == cty.NilType {
			v.Type = val.Type()
		}
	}

	// TODO krantzinator: not doing custom validations right now, unless it's easy
	// for _, block := range content.Blocks {
	// 	switch block.Type {

	// case "validation":
	// 	vv, moreDiags := decodeVariableValidationBlock(v.Name, block, override)
	// 	diags = append(diags, moreDiags...)
	// 	v.Validations = append(v.Validations, vv)

	// 	default:
	// 		// The above cases should be exhaustive for all block types
	// 		// defined in variableBlockSchema
	// 		panic(fmt.Sprintf("unhandled block type %q", block.Type))
	// 	}
	// }

	return &v, diags
}

// Prefix your environment variables with VarEnvPrefix so that Waypoint can see
// them.
const VarEnvPrefix = "WP_VAR_"

// SetJobInputValues collects values set via the CLI (-var, -varfile) and
// local env vars (WP_VAR_*) and translates those values to pb.Variables. These
// pb.Variables can then be set on the job for eventual parsing on the runner,
// after the runner has decoded the variables defined in the waypoint.hcl.
func SetJobInputValues(vars map[string]string, files []string) ([]*pb.Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	ret := []*pb.Variable{}

	{
		env := os.Environ()
		for _, raw := range env {
			if !strings.HasPrefix(raw, VarEnvPrefix) {
				continue
			}
			raw = raw[len(VarEnvPrefix):] // trim the prefix

			eq := strings.Index(raw, "=")
			if eq == -1 {
				// Seems invalid, so we'll ignore it.
				continue
			}

			name := raw[:eq]
			rawVal := raw[eq+1:]

			// TODO krantzinator: I'm making everything a Variable_Str which
			// removes the reason of even specifying a value type, so I should
			// not do that
			ret = append(ret, &pb.Variable{
				Name:   name,
				Value:  &pb.Variable_Str{Str: rawVal},
				Source: &pb.Variable_Env{},
			})
		}
	}

	for _, file := range files {
		if file != "" {
			pbv, diags := parseFileValues(file, sourceFile)
			if diags.HasErrors() {
				return nil, diags
			}
			ret = append(ret, pbv...)
		}
	}

	// Then we process values given explicitly on the command line, either
	// as individual literal settings or as files to read.
	for name, val := range vars {
		ret = append(ret, &pb.Variable{
			Name:   name,
			Value:  &pb.Variable_Str{Str: val},
			Source: &pb.Variable_Cli{},
		})
	}
	return ret, diags
}

// EvalInputValues evaluates the provided variable values and validates their
// types per the type declared in the waypoint.hcl for that variable name.
// The order in which values are evaluated corresponds to their precedence, with
// higher precedence values overwriting lower precedence values.
func EvalInputValues(
	pbvars []*pb.Variable,
	vs map[string]*Variable,
) (InputValues, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	iv := InputValues{}

	for v, def := range vs {
		iv[v] = def.Default
	}

	for _, pbv := range pbvars {
		variable, found := vs[pbv.Name]
		if !found {
			// TODO krantzinator: what to do with a warning diag type
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Undefined variable",
				Detail: fmt.Sprintf("A %q variable value was set, "+
					"but was not found in known variables. To declare variable "+
					"%q, place a variable definition block in your waypoint.hcl file.",
					pbv.Name, pbv.Name),
			})
			continue
		}

		var source string
		switch pbv.Source.(type) {
		case *pb.Variable_Cli:
			source = sourceCLI
		case *pb.Variable_File_:
			source = sourceFile
		case *pb.Variable_Env:
			source = sourceEnv
		case *pb.Variable_Vcs:
			source = sourceVCS
		case *pb.Variable_Server:
			source = sourceServer
		}

		var expr hclsyntax.Expression
		switch pbv.Value.(type) {

		// case *pb.Variable_Hcl:
		// 	value := pbv.Value.(*pb.Variable_Hcl).Hcl
		// 	fakeFilename := fmt.Sprintf("<value for var.%s from server>", pbv.Name)
		// 	expr, diags = hclsyntax.ParseExpression([]byte(value), fakeFilename, hcl.Pos{Line: 1, Column: 1})
		case *pb.Variable_Str:
			value := pbv.Value.(*pb.Variable_Str).Str
			expr = &hclsyntax.LiteralValueExpr{Val: cty.StringVal(value)}
		case *pb.Variable_Bool:
			value := pbv.Value.(*pb.Variable_Bool).Bool
			expr = &hclsyntax.LiteralValueExpr{Val: cty.BoolVal(value)}
		case *pb.Variable_Num:
			value := pbv.Value.(*pb.Variable_Num).Num
			expr = &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(value)}
		}

		val, valDiags := expr.Value(nil)
		diags = append(diags, valDiags...)
		if valDiags.HasErrors() {
			return nil, diags
		}

		if variable.Type != cty.NilType {
			var err error
			val, err = convert.Convert(val, variable.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid value for variable",
					Detail:   fmt.Sprintf("The value set for variable %q from source %q is not compatible with the variable's type constraint: %s.", pbv.Name, source, err),
				})
				val = cty.DynamicVal
			}
		}

		iv[pbv.Name] = &InputValue{
			Source: source,
			Value:  val,
			Expr:   expr,
		}
	}

	// check that all variables have a set value, including default of null
	for name, _ := range vs {
		_, found := iv[name]
		if !found {
			return nil, append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Unset variable %q", name),
				// TODO add our docs link here
				Detail: "A variable must be set or have a default value; see " +
					"[docs] for " +
					"details.",
			})
		}
	}

	return iv, diags
}

// LoadVCSFiles loads any *.auto.wpvars(.json) files in the VCS repo
func LoadVCSFiles(wd string) ([]*pb.Variable, hcl.Diagnostics) {
	var pbv []*pb.Variable
	var diags hcl.Diagnostics

	// Check working directory (vcs or local) for *.auto.wpvars(.json) files
	var varFiles []string
	if files, err := ioutil.ReadDir(wd); err == nil {
		for _, f := range files {
			name := f.Name()
			if !isAutoVarFile(name) {
				continue
			}
			varFiles = append(varFiles, wd+"/"+name)
		}
	}

	for _, f := range varFiles {
		if f != "" {
			pbv, diags = parseFileValues(f, sourceVCS)
			if diags.HasErrors() {
				return nil, diags
			}
		}
	}
	return pbv, nil
}

func parseFileValues(filename string, source string) ([]*pb.Variable, hcl.Diagnostics) {
	var pbv []*pb.Variable
	f, diags := readFileValues(filename)
	if diags.HasErrors() {
		return nil, diags
	}

	attrs, moreDiags := f.Body.JustAttributes()
	diags = append(diags, moreDiags...)
	// We grab all variables here; we'll later check set variables against the
	// known variables defined in the waypoint.hcl on the runner when we
	// consolidate values from local + server
	for name, attr := range attrs {
		val, moreDiags := attr.Expr.Value(nil)
		diags = append(diags, moreDiags...)

		v := &pb.Variable{
			Name: name,
		}
		// Set type
		switch val.Type() {
		case cty.String:
			v.Value = &pb.Variable_Str{Str: val.AsString()}
		case cty.Bool:
			v.Value = &pb.Variable_Bool{Bool: val.True()}
		case cty.Number:
			var num int64
			err := gocty.FromCtyValue(val, &num)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid number",
					Detail:   err.Error(),
					Subject:  &attr.Range,
				})
				return nil, diags
			}
			v.Value = &pb.Variable_Num{Num: num}
		}

		// Set source
		switch source {
		case sourceFile:
			v.Source = &pb.Variable_File_{}
		case sourceVCS:
			v.Source = &pb.Variable_Vcs{}
		}
		pbv = append(pbv, v)
	}

	return pbv, diags
}

// readFileValues is a helper function that loads a file, parses if it is
// hcl or json, and checks for any errant variable definition blocks. It returns
// the files contents for further evaluation if necessary.
func readFileValues(filename string) (*hcl.File, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	// load the file
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		var errStr string
		if os.IsNotExist(err) {
			errStr = fmt.Sprintf("Given variables file %s does not exist.", filename)
		} else {
			errStr = fmt.Sprintf("Error while reading %s: %s.", filename, err)
		}
		return nil, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read variable values from file",
			Detail:   errStr,
		})
	}

	// parse the file, whether it's hcl or json
	var f *hcl.File
	if strings.HasSuffix(filename, ".json") {
		var hclDiags hcl.Diagnostics
		f, hclDiags = hcljson.Parse(src, filename)
		diags = append(diags, hclDiags...)
		if f == nil || f.Body == nil {
			return nil, diags
		}
	} else {
		var hclDiags hcl.Diagnostics
		f, hclDiags = hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
		diags = append(diags, hclDiags...)
		if f == nil || f.Body == nil {
			return nil, diags
		}
	}

	// Before we do our real decode, we'll probe to see if there are any
	// blocks of type "variable" in this body, since it's a common mistake
	// for new users to put variable declarations in wpvars rather than
	// variable value definitions.
	{
		content, _, _ := f.Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type:       "variable",
					LabelNames: []string{"name"},
				},
			},
		})
		for _, block := range content.Blocks {
			name := block.Labels[0]
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Variable declaration in a wpvars file",
				Detail: fmt.Sprintf("A wpvars file is used to assign "+
					"values to variables that have already been declared "+
					"in the waypoint.hcl, not to declare new variables. To "+
					"declare variable %q, place this block in your "+
					"waypoint.hcl file.\n\nTo set a value for this variable "+
					"in %s, use the definition syntax instead:\n    %s = <value>",
					name, block.TypeRange.Filename, name),
				Subject: &block.TypeRange,
			})
		}
		if diags.HasErrors() {
			// If we already found problems then JustAttributes below will find
			// the same problems with less-helpful messages, so we'll bail for
			// now to let the user focus on the immediate problem.
			return nil, diags
		}
	}

	return f, diags
}

// // sortValues is a helper function that sorts variable Values by their
// // precedence order, ordering the slice from highest to lowest precedence
// func sortValues(values []InputValue) []InputValue {
// 	sort.Slice(values, func(i, j int) bool {
// 		return values[i].Source.Precedence > values[j].Source.Precedence
// 	})
// 	return values
// }

func (iv *InputValues) Values() map[string]cty.Value {
	res := map[string]cty.Value{}
	for k, v := range *iv {
		res[k] = v.Value
	}
	return res
}

func AddInputVariables(ctx *hcl.EvalContext, vs *InputValues) {
	vars := (*vs).Values()
	variables := map[string]cty.Value{
		"var": cty.ObjectVal(vars),
	}
	ctx.Variables = variables
}

// isAutoVarFile determines if the file ends with .auto.tfvars or .auto.tfvars.json
func isAutoVarFile(path string) bool {
	return strings.HasSuffix(path, ".auto.wpvars") ||
		strings.HasSuffix(path, ".auto.wpvars.json")
}
