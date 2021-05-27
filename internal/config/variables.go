package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// A consistent detail message for all "not a valid identifier" diagnostics.
const badIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."

// VariableAssignments contain the values from the job, from the server/VCS
// repo, and default values from the waypoint.hcl. These are to
// sort precedence and add the final variable values to the EvalContext
type VariableAssignment struct {
	Value  cty.Value
	Source string
	Expr   hcl.Expression
	// The location of the variable value if the value was provided
	// from a file
	Range  hcl.Range
}

// TODO krantzinator: this goes somewhere else
type Variable struct {
	Name string

	// Values contains possible values for the variable; The last value set
	// from these will be the one used. If none is set; an error will be
	// returned by Value().
	Values []VariableAssignment

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

type hclVariable struct {
	Name        string         `hcl:",label"`
	Default     cty.Value      `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
	Description string         `hcl:"description,optional"`
}

// Variables are used when parsing the Config, to set default values from
// the waypoint.hcl and bring in the values from the job and the server/VCS
// for eventual precedence sorting and setting on the EvalContext
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

func (c *Config) DecodeVariableBlocks(variables *Variables) hcl.Diagnostics {
	schema, _ := gohcl.ImpliedBodySchema(&validateStruct{})
	content, diag := c.hclConfig.Body.Content(schema)
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

// decodeVariableBlock decodes a "variable" block
// ectx is passed only in the evaluation of the default value.
func (variables *Variables) decodeVariableBlock(block *hcl.Block) hcl.Diagnostics {
	if (*variables) == nil {
		(*variables) = Variables{}
	}

	// TODO krantzinator: A lot of these validations happen twice before now --
	// when we validate the config on init and on the runner; we probably
	// don't need it here, too
	if _, found := (*variables)[block.Labels[0]]; found {
		return []*hcl.Diagnostic{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate variable",
			Detail:   "Duplicate " + block.Labels[0] + " variable definition found.",
			Context:  block.DefRange.Ptr(),
		}}
	}

	name := block.Labels[0]

	// TODO krantzinator this should happen earlier than the runner
	// Could be done before any ops are run and fail immediately
	content, diags := block.Body.Content(variableBlockSchema)
	if !hclsyntax.ValidIdentifier(name) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid variable name",
			Detail:   badIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	v := &Variable{
		Name:  name,
		Range: block.DefRange,
	}

	if attr, exists := content.Attributes["description"]; exists {
		valDiags := gohcl.DecodeExpression(attr.Expr, nil, &v.Description)
		diags = append(diags, valDiags...)
	}

	if t, ok := content.Attributes["type"]; ok {
		tp, moreDiags := typeexpr.Type(t.Expr)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return diags
		}

		v.Type = tp
	}

	if def, ok := content.Attributes["default"]; ok {
		defaultValue, moreDiags := def.Expr.Value(nil)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return diags
		}

		if v.Type != cty.NilType {
			var err error
			defaultValue, err = convert.Convert(defaultValue, v.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("This default value is not compatible with the variable's type constraint: %s.", err),
					Subject:  def.Expr.Range().Ptr(),
				})
				defaultValue = cty.DynamicVal
			}
		}

		v.Values = append(v.Values, VariableAssignment{
			Source: "default",
			Value:  defaultValue,
		})

		// It's possible no type attribute was assigned so lets make sure we
		// have a valid type otherwise there could be issues parsing the value.
		if v.Type == cty.NilType {
			v.Type = defaultValue.Type()
		}
	}

	(*variables)[name] = v

	return diags
}

// Prefix your environment variables with VarEnvPrefix so that Waypoint can see
// them.
const VarEnvPrefix = "WP_VAR_"

// Collect values from flags, local var files, and env vars
func CollectInputVars(vars map[string]string, files []string) ([]*pb.Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	ret := []*pb.Variable{}

	// First we'll deal with environment variables, since they have the lowest
	// precedence.
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
	// We'll first check for a value already set and overwrite if there, or
	// append a new value
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
func (variables *Variables) CollectInputValRemote(files []*hcl.File, serverVars []*pb.Variable) hcl.Diagnostics {
	var diags hcl.Diagnostics

	// files will contain files found in the remote git source
	for _, file := range files {
		// Before we do our real decode, we'll probe to see if there are any
		// blocks of type "variable" in this body, since it's a common mistake
		// for new users to put variable declarations in wpvars rather than
		// variable value definitions.
		{
			content, _, _ := file.Body.PartialContent(&hcl.BodySchema{
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
					Summary:  "Variable declaration in a .wpvars file",
					Detail: fmt.Sprintf("A .wpvars file is used to assign "+
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

		attrs, moreDiags := file.Body.JustAttributes()
		diags = append(diags, moreDiags...)

		for name, attr := range attrs {
			variable, found := (*variables)[name]
			if !found {
				sev := hcl.DiagWarning
				// TODO krantzinator
				// if cfg.ValidationOptions.Strict {
				// 	sev = hcl.DiagError
				// }
				diags = append(diags, &hcl.Diagnostic{
					Severity: sev,
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

			variable.Values = append(variable.Values, VariableAssignment{
				Source: "repofile",
				Value:  val,
				Expr:   attr.Expr,
			})
		}
	}

	// Finally we process values given explicitly on the command line.
	for _, sv := range serverVars {
		variable, found := (*variables)[sv.Name]
		if !found {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Undefined server variable",
				Detail: fmt.Sprintf("A %q variable was set in the UI "+
					"and stored on the server, but was not found in "+
					"known variables. To declare variable %q, place "+
					"this block in your waypoint.hcl file.",
					sv.Name, sv.Name),
			})
			continue
		}

		var expr hclsyntax.Expression
		switch sv.Value.(type) {

		case *pb.Variable_Hcl:
			value := sv.Value.(*pb.Variable_Hcl).Hcl
			fakeFilename := fmt.Sprintf("<value for var.%s from server>", sv.Name)
			expr, diags = hclsyntax.ParseExpression([]byte(value), fakeFilename, hcl.Pos{Line: 1, Column: 1})

		case *pb.Variable_Str:
			value := sv.Value.(*pb.Variable_Str).Str
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
					Summary:  "Invalid argument value for server-stored variable",
					Detail:   fmt.Sprintf("The received arg value for %s is not compatible with the variable's type constraint: %s.", sv.Name, err),
					Subject:  expr.Range().Ptr(),
				})
				val = cty.DynamicVal
			}
		}

		variable.Values = append(variable.Values, VariableAssignment{
			Source: "server",
			Value:  val,
			Expr:   expr,
		})
	}

	return diags
}

func (variables *Variables) SortPrecedence() ([]*pb.Variable, error) {
	var ret []*pb.Variable

	for _, v := range *variables {
		for _, vv := range v.Values {
			switch st := vv.Source; st {
			case "default":
				pbv := &pb.Variable{
					Name:   v.Name,
					Value:  &pb.Variable_Str{Str: vv.Value.AsString()},
					Source: nil,
				}
				ret = append(ret, pbv)
			case "server":
				pbv := &pb.Variable{
					Name:   v.Name,
					Value:  &pb.Variable_Str{Str: vv.Value.AsString()},
					Source: &pb.Variable_Server{},
				}
				ret = append(ret, pbv)
			}
		}
	}

	return ret, nil
}
