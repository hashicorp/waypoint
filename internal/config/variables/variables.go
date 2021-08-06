package variables

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"
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

const (
	// Prefix for collecting variable values from environment variables
	varEnvPrefix = "WP_VAR_"

	// Variable value sources
	// listed in descending precedence order for ease of reference
	sourceCLI     = "cli"
	sourceFile    = "file"
	sourceEnv     = "env"
	sourceVCS     = "vcs"
	sourceServer  = "server"
	sourceDefault = "default"
)

var (
	// sourceMap maps a variable pb source type to its string representation
	fromSource = map[reflect.Type]string{
		reflect.TypeOf((*pb.Variable_Cli)(nil)):    sourceCLI,
		reflect.TypeOf((*pb.Variable_File_)(nil)):  sourceFile,
		reflect.TypeOf((*pb.Variable_Env)(nil)):    sourceEnv,
		reflect.TypeOf((*pb.Variable_Vcs)(nil)):    sourceVCS,
		reflect.TypeOf((*pb.Variable_Server)(nil)): sourceServer,
	}

	// The attributes we expect to see in variable blocks
	// Future expansion here could include `sensitive`, `validations`, etc
	variableBlockSchema = &hcl.BodySchema{
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
)

// Variable stores a parsed variable definition from the waypoint.hcl
type Variable struct {
	Name string

	// The default value in the variable definition
	Default *Value

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

// HclVariable is used when decoding the waypoint.hcl config. Because we use
// hclsimple for this decode, we need the `Type` to be evaluated as an hcl
// expression. When we parse the config, we need `Type` to be evaluated as
// cty.Type, so this struct is only used for the basic decoding of the file
// to verify HCL syntax.
type HclVariable struct {
	Name        string         `hcl:",label"`
	Default     cty.Value      `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
	Description string         `hcl:"description,optional"`
}

// Values are used to store values collected from various sources.
// Values are added to the map in precedence order, and then used to
// create the final map of cty.Values for config hcl context evaluation.
type Values map[string]*Value

// Value contain the value of the variable along with associated metada,
// including the source it was set from: cli, file, env, vcs, server/ui
type Value struct {
	Value  cty.Value
	Source string
	Expr   hcl.Expression
	// The location of the variable value if the value was provided from a file
	Range hcl.Range
}

// DecodeVariableBlocks uses the hclConfig schema to iterate over all
// variable blocks, validating names and types and checking for duplicates.
// It returns the final map of Variables to store for later reference.
func DecodeVariableBlocks(content *hcl.BodyContent) (map[string]*Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	vs := map[string]*Variable{}
	for _, block := range content.Blocks.OfType("variable") {
		v, diags := decodeVariableBlock(block)
		if diags.HasErrors() {
			return nil, diags
		}

		if _, found := vs[v.Name]; found {
			return nil, []*hcl.Diagnostic{{
				Severity: hcl.DiagError,
				Summary:  "Duplicate variable",
				Detail:   "Duplicate " + v.Name + " variable definition found.",
				Context:  block.DefRange.Ptr(),
			}}
		}

		vs[block.Labels[0]] = v
	}

	return vs, diags
}

// decodeVariableBlock validates each part of the variable block,
// building out a defined *Variable
func decodeVariableBlock(block *hcl.Block) (*Variable, hcl.Diagnostics) {
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
			Detail:   "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes.",
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

		v.Default = &Value{
			Source: sourceDefault,
			Value:  val,
		}

		// It's possible no type attribute was assigned so lets make sure we
		// have a valid type otherwise there could be issues parsing the value.
		if v.Type == cty.NilType {
			v.Type = val.Type()
		}
	}

	return &v, diags
}

// LoadVariableValues collects values set via the CLI (-var, -varfile) and
// local env vars (WP_VAR_*) and translates those values to pb.Variables. These
// pb.Variables can then be set on the job for eventual parsing on the runner,
// after the runner has decoded the variables defined in the waypoint.hcl.
// All values are set as protobuf strings, with the expectation that later
// evaluation will convert them to their defined types.
func LoadVariableValues(vars map[string]string, files []string) ([]*pb.Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	ret := []*pb.Variable{}

	// The order here is important, as the order in which values are evaluated
	// dictate their precedence. We therefore evalute these three sources in order
	// of env, file, cli.

	// process env values ("env" source)
	{
		env := os.Environ()
		for _, raw := range env {
			if !strings.HasPrefix(raw, varEnvPrefix) {
				continue
			}
			raw = raw[len(varEnvPrefix):] // trim the prefix

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

	// process -var-file args ("file" source)
	for _, file := range files {
		if file != "" {
			pbv, diags := parseFileValues(file, sourceFile)
			if diags.HasErrors() {
				return nil, diags
			}
			ret = append(ret, pbv...)
		}
	}

	// process -var args ("cli" source)
	for name, val := range vars {
		ret = append(ret, &pb.Variable{
			Name:   name,
			Value:  &pb.Variable_Str{Str: val},
			Source: &pb.Variable_Cli{},
		})
	}
	return ret, diags
}

// EvaluateVariables evaluates the provided variable values and validates their
// types per the type declared in the waypoint.hcl for that variable name.
// The order in which values are evaluated corresponds to their precedence, with
// higher precedence values overwriting lower precedence values.
// The supplied map of *Variable should be all defined variables (currently
// comes from decoding all variable blocks within the waypoint.hcl), and
// is used to validate types and that all variables have at least one
// assigned value.
func EvaluateVariables(
	pbvars []*pb.Variable,
	vs map[string]*Variable,
	log hclog.Logger,
) (Values, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	iv := Values{}

	for v, def := range vs {
		iv[v] = def.Default
	}

	for _, pbv := range pbvars {
		variable, found := vs[pbv.Name]
		if !found {
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

		// set our source for error messaging
		source := fromSource[reflect.TypeOf(pbv.Source)]
		if source == "" {
			source = "unknown"
			log.Debug("No source found for value given for variable %q", pbv.Name)
		}

		// We have to specify the three different simple types we support -- string,
		// bool, number -- when doing the below evaluation of hcl expressions
		// because of our translation to-and-from protobuf format.
		// While cty allows us to parse all simple types as LiteralValueExpr, we
		// have to first translate the pb values back into cty values, thus
		// necessitating a separate case statement for each simple type
		var expr hclsyntax.Expression
		switch sv := pbv.Value.(type) {
		case *pb.Variable_Hcl:
			value := sv.Hcl
			fakeFilename := fmt.Sprintf("<value for var.%s from source %q>", pbv.Name, source)
			expr, diags = hclsyntax.ParseExpression([]byte(value), fakeFilename, hcl.Pos{Line: 1, Column: 1})
		case *pb.Variable_Str:
			value := sv.Str
			expr = &hclsyntax.LiteralValueExpr{Val: cty.StringVal(value)}
		case *pb.Variable_Bool:
			value := sv.Bool
			expr = &hclsyntax.LiteralValueExpr{Val: cty.BoolVal(value)}
		case *pb.Variable_Num:
			value := sv.Num
			expr = &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(value)}
		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid value type for variable",
				Detail:   "The variable type was not set as a string, number, bool, or hcl expression",
			})
			return nil, diags
		}

		val, valDiags := expr.Value(nil)
		if valDiags.HasErrors() {
			diags = append(diags, valDiags...)
			return nil, diags
		}

		if variable.Type != cty.NilType {
			var err error
			// store the current cty.Value type before attempting the conversion
			valType := val.Type().FriendlyName()

			// If the value came from the cli or an env var, it was stored as
			// a raw string value; however it could be a complex value, such as a
			// map/list/etc.
			// Now that we know the expected type, we'll check here for that
			// and, if necessary, repeat the expression parsing for HCL syntax
			if source == sourceCLI || source == sourceEnv {
				if !variable.Type.IsPrimitiveType() {
					fakeFilename := fmt.Sprintf("<value for var.%s from source %q>", pbv.Name, source)
					expr, diags = hclsyntax.ParseExpression([]byte(val.AsString()), fakeFilename, hcl.Pos{Line: 1, Column: 1})
					val, valDiags = expr.Value(nil)
					if valDiags.HasErrors() {
						diags = append(diags, valDiags...)
						return nil, diags
					}
				}
			}

			val, err = convert.Convert(val, variable.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid value for variable",
					Detail: fmt.Sprintf(
						"The value set for variable %q from source %q is of value type %q and is not compatible with the variable's type constraint: %s.",
						pbv.Name,
						source,
						valType,
						err,
					),
				})
				val = cty.DynamicVal
			}
		}

		iv[pbv.Name] = &Value{
			Source: source,
			Value:  val,
			Expr:   expr,
		}
	}

	// check that all variables have a set value, including default of null
	for name := range vs {
		v, ok := iv[name]
		if !ok || v == nil {
			return nil, append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Unset variable %q", name),
				// TODO krantzinator: add our docs link here
				Detail: "A variable must be set or have a default value; see " +
					"[docs] for " +
					"details.",
			})
		}
	}

	return iv, diags
}

// LoadAutoFiles loads any *.auto.wpvars(.json) files in the source repo
func LoadAutoFiles(wd string) ([]*pb.Variable, hcl.Diagnostics) {
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
			varFiles = append(varFiles, filepath.Join(wd, name))
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

// parseFileValues is a helper function to extract variable values from the
// provided file, using the provided source to set the pb.Variable.Source value.
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
		default:
			// if it's not a primitive/simple type, we set as hcl here to be later
			// evaluated as hcl; any errors at evaluating the hcl type will
			// be handled at that time
			v.Value = &pb.Variable_Hcl{Hcl: val.AsString()}
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
// the files contents for further evaluation.
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
			return nil, diags
		}
	}
	return f, diags
}

// values creates a map of cty.values from the map of InputValues, for use
// in creating hcl contexts
func (iv Values) values() map[string]cty.Value {
	res := map[string]cty.Value{}
	for k, v := range iv {
		res[k] = v.Value
	}
	return res
}

// AddInputVariables adds the InputValues to the provided hcl context
func AddInputVariables(ctx *hcl.EvalContext, vs Values) {
	vars := vs.values()
	variables := map[string]cty.Value{
		"var": cty.ObjectVal(vars),
	}
	ctx.Variables = variables
}

// isAutoVarFile determines if the file ends with .auto.wpvars or .auto.wpvars.json
func isAutoVarFile(path string) bool {
	return strings.HasSuffix(path, ".auto.wpvars") ||
		strings.HasSuffix(path, ".auto.wpvars.json")
}
