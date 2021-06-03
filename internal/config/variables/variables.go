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
)

// A consistent detail message for all "not a valid identifier" diagnostics.
const (
	badIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."

	// Variable value sources
	sourceDefault = "default"
	sourceCLI     = "cli"
	sourceFile    = "file"
	sourceEnv     = "env"
	sourceVCS     = "vcs"
	sourceServer  = "server"
)

// Value contain the value of the variable along with associated metada,
// including the source it was set from: cli, file, env, vcs, server/ui
type Value struct {
	Value  cty.Value
	Source string
	Expr   hcl.Expression
	// The location of the variable value if the value was provided
	// from a file
	Range hcl.Range
}

type Variable struct {
	// Value contain the values from the job, from the server/VCS
	// repo, and default values from the waypoint.hcl. These are to
	// sort precedence and add the final variable values to the EvalContext
	Values []Value

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

type HclBase struct {
	Variables []*HclVariable `hcl:"variable,block"`
	Body      hcl.Body       `hcl:",body"`
	Remain    hcl.Body       `hcl:",remain"`
}

type HclVariable struct {
	Name        string         `hcl:",label"`
	Default     cty.Value      `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
	Description string         `hcl:"description,optional"`
}

// Variables are used when parsing the Config, to set default values from
// the waypoint.hcl and bring in the values from the job and the server/VCS
// for eventual precedence sorting and setting on the EvalContext
// TODO krantzinator: make these InputVars
type Variables map[string]*Variable

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

func (variables *Variables) DecodeVariableBlocks(body hcl.Body) hcl.Diagnostics {
	schema, _ := gohcl.ImpliedBodySchema(&HclBase{})
	content, diag := body.Content(schema)
	if diag.HasErrors() {
		return diag
	}

	var diags hcl.Diagnostics
	for _, block := range content.Blocks.OfType("variable") {
		moreDiags := variables.decodeVariableBlock(block)
		if moreDiags != nil {
			diags = append(diags, moreDiags...)
		}
	}

	return diags
}

// decodeVariableBlock first validates the variable block, and then sets
// the decoded variables with their default value on *Variables
func (variables *Variables) decodeVariableBlock(block *hcl.Block) hcl.Diagnostics {
	if (*variables) == nil {
		(*variables) = Variables{}
	}
	name := block.Labels[0]

	// Checking for duplicates happens here, rather than during the config.Validate
	// step, because config.Validate doesn't store any decoded blocks.
	if _, found := (*variables)[name]; found {
		return []*hcl.Diagnostic{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate variable",
			Detail:   "Duplicate " + name + " variable definition found.",
			Context:  block.DefRange.Ptr(),
		}}
	}

	v, diags := ValidateVarBlock(block)
	if diags.HasErrors() {
		return diags
	}

	(*variables)[name] = v

	return diags
}

// ValidateVarBlock validates each part of the variable block, building out
// the final *Variable along the way
func ValidateVarBlock(block *hcl.Block) (*Variable, hcl.Diagnostics) {
	v := Variable{
		Range: block.DefRange,
	}

	content, diags := block.Body.Content(variableBlockSchema)

	if !hclsyntax.ValidIdentifier(block.Labels[0]) {
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

	if t, ok := content.Attributes["type"]; ok {
		tt, moreDiags := typeexpr.Type(t.Expr)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags
		}
		v.Type = tt
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
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("This default value is not compatible with the variable's type constraint: %s.", err),
					Subject:  attr.Expr.Range().Ptr(),
				})
				return nil, diags
			}
			val = cty.DynamicVal
		}

		// TODO krantzinator; may not do this slice of var assignments
		v.Values = append(v.Values, Value{
			Source: sourceDefault,
			Value:  val,
		})

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

// Collect values from flags, local var files, and env vars
func CollectInputVars(vars map[string]string, files []string) ([]*pb.Variable, hcl.Diagnostics) {
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

			ret = append(ret, &pb.Variable{
				Name:   name,
				Value:  &pb.Variable_Str{Str: rawVal},
				Source: &pb.Variable_Env{},
			})
		}
	}

	// TODO var file parsing

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

// Collect values from server and remote-stored wpvars file
func (variables *Variables) CollectInputValues(files []string, pbvars []*pb.Variable) hcl.Diagnostics {
	var diags hcl.Diagnostics

	// files will contain files found in the remote git source
	for _, file := range files {
		if file != "" {
			fileDiags := variables.parseFileValues(file)
			diags = append(diags, fileDiags...)
			if diags.HasErrors() {
				return diags
			}
		}
	}

	for _, pbv := range pbvars {
		variable, found := (*variables)[pbv.Name]
		if !found {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Undefined variable",
				Detail: fmt.Sprintf("A %q variable value was set, "+
					"but was not found in known variables. To declare "+
					"variable %q, place this block in your waypoint.hcl file.",
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

		case *pb.Variable_Hcl:
			value := pbv.Value.(*pb.Variable_Hcl).Hcl
			fakeFilename := fmt.Sprintf("<value for var.%s from server>", pbv.Name)
			expr, diags = hclsyntax.ParseExpression([]byte(value), fakeFilename, hcl.Pos{Line: 1, Column: 1})

		case *pb.Variable_Str:
			value := pbv.Value.(*pb.Variable_Str).Str
			expr = &hclsyntax.LiteralValueExpr{Val: cty.StringVal(value)}
		}

		val, valDiags := expr.Value(nil)
		diags = append(diags, valDiags...)
		if valDiags.HasErrors() {
			return diags
		}

		if variable.Type != cty.NilType {
			var err error
			val, err = convert.Convert(val, variable.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid value for variable",
					Detail:   fmt.Sprintf("The value set for variable %q from source %q is not compatible with the variable's type constraint: %s.", pbv.Name, source, err),
					Subject:  expr.Range().Ptr(),
				})
				val = cty.DynamicVal
			}
		}

		variable.Values = append(variable.Values, Value{
			Source: source,
			Value:  val,
			Expr:   expr,
		})
	}

	return diags
}

func (variables *Variables) SortPrecedence(vars []*pb.Variable) error {

	//

	// for _, v := range *variables {
	// 	// gsv.Name
	// }

	return nil
}

func (variables *Variables) parseFileValues(filename string) hcl.Diagnostics {
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
		return append(diags, &hcl.Diagnostic{
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
			return diags
		}
	} else {
		var hclDiags hcl.Diagnostics
		f, hclDiags = hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
		diags = append(diags, hclDiags...)
		if f == nil || f.Body == nil {
			return diags
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
			return diags
		}
	}

	attrs, moreDiags := f.Body.JustAttributes()
	diags = append(diags, moreDiags...)

	for name, attr := range attrs {
		variable, found := (*variables)[name]
		if !found {
			// TODO krantzinator: what to do with a warning diag type
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Undefined variable",
				Detail: fmt.Sprintf("A %q variable was set but was "+
					"not found in known variables. To declare "+
					"variable %q, place this block in your "+
					"waypoint.hcl file.",
					name, name),
				Context: attr.Range.Ptr(),
			})
			continue
		}

		val, moreDiags := attr.Expr.Value(nil)
		diags = append(diags, moreDiags...)

		if variable.Type != cty.NilType {
			var err error
			val, err = convert.Convert(val, variable.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid value for variable",
					Detail:   fmt.Sprintf("The value for %s is not compatible with the variable's type constraint: %s.", name, err),
					Subject:  attr.Expr.Range().Ptr(),
				})
				val = cty.DynamicVal
			}
		}

		variable.Values = append(variable.Values, Value{
			Source: sourceFile,
			Value:  val,
			Expr:   attr.Expr,
		})
	}

	return diags
}
