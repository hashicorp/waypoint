// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package variables

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/hashicorp/waypoint/internal/appconfig"
	"github.com/hashicorp/waypoint/internal/config/dynamic"
	"github.com/hashicorp/waypoint/internal/config/variables/formatter"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	// Prefix for collecting variable values from environment variables
	varEnvPrefix = "WP_VAR_"
)

var (
	// sourceMap maps a variable pb source type to its string representation
	fromSource = map[reflect.Type]string{
		reflect.TypeOf((*pb.Variable_Cli)(nil)):     formatter.SourceCLI,
		reflect.TypeOf((*pb.Variable_File_)(nil)):   formatter.SourceFile,
		reflect.TypeOf((*pb.Variable_Env)(nil)):     formatter.SourceEnv,
		reflect.TypeOf((*pb.Variable_Vcs)(nil)):     formatter.SourceVCS,
		reflect.TypeOf((*pb.Variable_Server)(nil)):  formatter.SourceServer,
		reflect.TypeOf((*pb.Variable_Dynamic)(nil)): formatter.SourceDynamic,
	}

	fromSourceToFV = map[string]pb.Variable_FinalValue_Source{
		formatter.SourceCLI:     pb.Variable_FinalValue_CLI,
		formatter.SourceFile:    pb.Variable_FinalValue_FILE,
		formatter.SourceEnv:     pb.Variable_FinalValue_ENV,
		formatter.SourceVCS:     pb.Variable_FinalValue_VCS,
		formatter.SourceServer:  pb.Variable_FinalValue_SERVER,
		formatter.SourceDynamic: pb.Variable_FinalValue_DYNAMIC,
		formatter.SourceDefault: pb.Variable_FinalValue_DEFAULT,
		formatter.SourceUnknown: pb.Variable_FinalValue_UNKNOWN,
	}

	// The attributes we expect to see in variable blocks
	// Future expansion here could include `validations`, etc
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
			{
				Name: "env",
			},
			{
				Name: "sensitive",
			},
		},
	}
)

// Variable stores a parsed variable definition from the waypoint.hcl
type Variable struct {
	Name string

	// The default value in the variable definition
	Default *Value

	// A list of environment variables that will be sourced to satisfy
	// the value of this variable.
	Env []string

	// Cty Type of the variable. If the default value or a collected value is
	// not of this type nor can be converted to this type an error diagnostic
	// will show up. This allows us to assume that values are valid.
	//
	// When a default value - and no type - is passed into the variable
	// declaration, the type of the default variable will be used.
	Type cty.Type

	// Variables with this set will be hashed as SHA256 values for
	// the purposes of output and logging
	Sensitive bool

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
	Default     hcl.Expression `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
	Description string         `hcl:"description,optional"`
	Env         []string       `hcl:"env,optional"`
	Sensitive   bool           `hcl:"sensitive,optional"`
}

// Values are used to store values collected from various sources.
// Values are added to the map in precedence order, and then used to
// create the final map of cty.Values for config hcl context evaluation.
type Values map[string]*Value

// Value contain the value of the variable along with associated metadata,
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
func DecodeVariableBlocks(
	ctx *hcl.EvalContext,
	content *hcl.BodyContent,
) (map[string]*Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	vs := map[string]*Variable{}
	for _, block := range content.Blocks.OfType("variable") {
		v, diags := decodeVariableBlock(ctx, block)
		if diags.HasErrors() {
			return nil, diags
		}

		if _, found := vs[v.Name]; found {
			return nil, []*hcl.Diagnostic{{
				Severity: hcl.DiagError,
				Summary:  "Duplicate variable",
				Detail:   "Duplicate " + v.Name + " variable definition found.",
				Subject:  &v.Range,
				Context:  block.DefRange.Ptr(),
			}}
		}

		vs[block.Labels[0]] = v
	}

	return vs, diags
}

// decodeVariableBlock validates each part of the variable block,
// building out a defined *Variable
func decodeVariableBlock(
	ctx *hcl.EvalContext,
	block *hcl.Block,
) (*Variable, hcl.Diagnostics) {
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
		// TypeConstraint allows "any", and it's OK if users opt out of
		// waypoint type checking here.
		t, moreDiags := typeexpr.TypeConstraint(attr.Expr)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags
		}
		v.Type = t
	}

	if attr, exists := content.Attributes["env"]; exists {
		valDiags := gohcl.DecodeExpression(attr.Expr, nil, &v.Env)
		diags = append(diags, valDiags...)
		if diags.HasErrors() {
			return nil, diags
		}
	}

	if attr, exists := content.Attributes["sensitive"]; exists {
		valDiags := gohcl.DecodeExpression(attr.Expr, nil, &v.Sensitive)
		diags = append(diags, valDiags...)
	}

	if attr, exists := content.Attributes["default"]; exists {
		defaultCtx := ctx.NewChild()
		defaultCtx.Functions = dynamic.Register(map[string]function.Function{})

		val, valDiags := attr.Expr.Value(defaultCtx)
		diags = append(diags, valDiags...)
		if diags.HasErrors() {
			return nil, diags
		}

		// Depending on the value type, we behave differently.
		switch val.Type() {
		case dynamic.Type:
			// Dynamic types can either be strings, or values that can be
			// represented as json (i.e. objects and maps). The only allowed
			// primitive type is string, and users can use explicit
			// type conversion such as `tonumber` to achieve other primitive types.

			// This is intended to help users catch invalid types early. If
			// the type they specified doesn't match the actual type at runtime,
			// they'll get another error and that's OK.
			switch v.Type {
			case cty.Number, cty.Bool:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("type for variable %q must be string for dynamic values", name),
					Detail: "When using dynamically sourced configuration values, " +
						"can either be strings, or complex types. If you're unsure of " +
						"which, consult your configsourcer plugin documentation. If you want " +
						"to represent a string as another kind of value, use type " +
						"conversion functions such as `tonumber` when using the variable.",
					Subject: attr.Expr.Range().Ptr(),
				})
				val = cty.DynamicVal
			}

		default:
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
		}

		v.Default = &Value{
			Source: formatter.SourceDefault,
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
			pbv, diags := parseFileValues(file, formatter.SourceFile)
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

// LoadEnvValues loads the variable values from environment variables
// specified via the `env` field on the `variable` stanza.
func LoadEnvValues(vars map[string]*Variable) ([]*pb.Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	var ret []*pb.Variable

	for _, variable := range vars {
		// First we check for the WP_VAR_ value cause that always wins.
		v := os.Getenv(varEnvPrefix + variable.Name)

		// If we didn't find one and we have other sources, check those.
		if v == "" && len(variable.Env) > 0 {
			for _, env := range variable.Env {
				v = os.Getenv(env)
				if v != "" {
					break
				}
			}
		}

		// If we still have no value, then we set nothing.
		if v == "" {
			continue
		}

		ret = append(ret, &pb.Variable{
			Name:   variable.Name,
			Value:  &pb.Variable_Str{Str: v},
			Source: &pb.Variable_Env{},
		})
	}

	return ret, diags
}

// NeedsDynamicDefaults returns true if there are variables with a dynamic
// default value set that must be evaluated (because the value is not
// overridden).
func NeedsDynamicDefaults(
	pbvars []*pb.Variable,
	vars map[string]*Variable,
) bool {
	// Get all our variables with dynamic defaults
	dynamicVars := map[string]*Variable{}
	for k, v := range vars {
		if v.Default != nil {
			val := v.Default.Value
			if val.Type() == dynamic.Type {
				dynamicVars[k] = v
			}
		}
	}

	// Go through our variable values and delete any dynamic vars we have
	// values for already; we do not need to fetch those.
	for _, pbv := range pbvars {
		delete(dynamicVars, pbv.Name)
	}

	return len(dynamicVars) > 0
}

// LoadDynamicDefaults will load the default values for variables that have
// dynamic configurations. This will only load the values if there isn't an
// existing variable set in pbvars. Therefore, it is recommended that this is
// called last and the values are _prepended_ to pbvars for priority.
func LoadDynamicDefaults(
	ctx context.Context,
	log hclog.Logger,
	pbvars []*pb.Variable,
	cfgSrcs []*pb.ConfigSource,
	vars map[string]*Variable,
	dynamicOpts ...appconfig.Option,
) ([]*pb.Variable, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	// Get all our variables with dynamic defaults
	dynamicVars := map[string]*Variable{}
	for k, v := range vars {
		if v.Default != nil {
			val := v.Default.Value
			if val.Type() == dynamic.Type {
				dynamicVars[k] = v
			}
		}
	}

	// If we have no dynamic vars, we do nothing.
	if len(dynamicVars) == 0 {
		log.Debug("no dynamic vars")
		return nil, diags
	}

	log.Debug("dynamic variables discovered", "total", len(dynamicVars))

	// Go through our variable values and delete any dynamic vars we have
	// values for already; we do not need to fetch those.
	for _, pbv := range pbvars {
		delete(dynamicVars, pbv.Name)
	}

	// If we have no dynamic vars we need values for, also do nothing.
	if len(dynamicVars) == 0 {
		log.Debug("no dynamic vars needed values")
		return nil, diags
	}

	// Build our watcher
	ch := make(chan *appconfig.UpdatedConfig)
	w, err := appconfig.NewWatcher(append(
		append([]appconfig.Option{}, dynamicOpts...),
		appconfig.WithLogger(log),
		appconfig.WithNotify(ch),
	)...)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Error initializing configuration watcher",
			Detail:   err.Error(),
			Subject: &hcl.Range{
				Filename: "waypoint.hcl",
			},
		})
		return nil, diags
	}
	defer w.Close()

	// We have some variables we need to fetch dynamic values for.
	configVars := make([]*pb.ConfigVar, 0, len(dynamicVars))
	for k, v := range dynamicVars {
		configVars = append(configVars, &pb.ConfigVar{
			Name: k,
			Value: &pb.ConfigVar_Dynamic{
				Dynamic: v.Default.Value.EncapsulatedValue().(*pb.ConfigVar_DynamicVal),
			},

			// This is set to true on purpose because it forces the appconfig
			// watcher to give us an easier to consume format (struct vs
			// array of key=value strings for env vars).
			NameIsPath: true,
		})
	}

	// Update and send any config source overrides for dynamic vars.
	w.UpdateSources(ctx, cfgSrcs)

	// Send our variables. Purposely ignore the error return value because
	// it can only ever be a context cancellation which we pick up in the
	// select later.
	w.UpdateVars(ctx, configVars)

	// Wait for values.
	log.Debug("waiting for dynamic variable values", "count", len(configVars))
	for {
		select {
		case <-ctx.Done():
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Cancellation while waiting for configuration.",
				Detail:   ctx.Err().Error(),
				Subject: &hcl.Range{
					Filename: "waypoint.hcl",
				},
			})
			return nil, diags

		case <-time.After(5 * time.Second):
			log.Warn("waiting for dynamic variables, this delay is usually due to external systems")

		case config := <-ch:
			var result []*pb.Variable

			for _, f := range config.Files {

				if matchingVar, ok := vars[f.Path]; ok {
					if !matchingVar.Type.Equals(cty.String) {

						// This is some hacking. The user put something other than
						// `type = string` as the type for this variable.
						// The plugin has done one of three things:
						//
						// 1: Returned json typed data (e.g. plugin.proto ConfigSource -> Value -> Value -> result -> json)
						//    This is the happy path.
						//    JSON is valid HCL, so we can run it up as an hcl variable,
						//    it'll get properly marshalled and type checked, and can be
						//    used elsewhere in the waypoint hcl.
						// 2: The plugin returned a string value, but the string is
						//    actually json. This will also work, as long as the json
						//    and user-specified types match.
						// 3: The plugin returned a non-json string. In this case, we
						//    emit a runtime error.

						// Someday, it would be nice to KNOW here if the plugin returned structured
						// data or not. That's work though, and option 2 is a happy side effect
						// of us not knowing, so it's OK for today.

						// Make sure the plugin returned valid non-string data.
						if err := json.Unmarshal(f.Data, &json.RawMessage{}); err != nil {
							diags = append(diags, &hcl.Diagnostic{
								Severity: hcl.DiagError,
								Summary:  "Plugin output <> hcl type mismatch",
								Detail: fmt.Sprintf(
									"Variable %q is declared as non-string type %q, \n"+
										"but the configsourcer plugin did not return\n"+
										"structured data that can be json-marshalled:\n%s",
									matchingVar.Name,
									matchingVar.Type.FriendlyNameForConstraint(),
									err,
								),
								Subject: &hcl.Range{
									Filename: "waypoint.hcl",
								},
							})
							return nil, diags
						}

						result = append(result, &pb.Variable{
							Name: f.Path,
							Value: &pb.Variable_Hcl{
								// json is valid hcl!
								Hcl: string(f.Data),
							},
							Source: &pb.Variable_Dynamic{},
						})
					} else {
						result = append(result, &pb.Variable{
							Name: f.Path,
							Value: &pb.Variable_Str{
								Str: string(f.Data),
							},
							Source: &pb.Variable_Dynamic{},
						})
					}
				}
			}

			return result, diags
		}
	}
}

// EvaluateVariables evaluates the provided variable values and validates their
// types per the type declared in the waypoint.hcl for that variable name.
// The order in which values are evaluated corresponds to their precedence, with
// higher precedence values overwriting lower precedence values.
//
// The supplied map of *Variable should be all defined variables (currently
// comes from decoding all variable blocks within the waypoint.hcl), and
// is used to validate types and that all variables have at least one
// assigned value.
func EvaluateVariables(
	log hclog.Logger,
	pbvars []*pb.Variable,
	vs map[string]*Variable,
	salt string,
) (Values, map[string]*pb.Variable_FinalValue, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	iv := Values{}

	for v, def := range vs {
		// Do not allow dynamic values as default values since they aren't valid.
		// Dynamic values should be evaluated and overridden by LoadDynamicDefaults
		// and provided via pbvars. If not, then an unset error will be created.
		if def.Default != nil && def.Default.Value.Type() == dynamic.Type {
			continue
		}

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
				Subject: &hcl.Range{
					Filename: "waypoint.hcl",
				},
			})
			continue
		}

		// set our source for error messaging
		source := fromSource[reflect.TypeOf(pbv.Source)]
		if source == "" {
			source = formatter.SourceUnknown
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
				Subject:  &variable.Range,
			})
			return nil, nil, diags
		}

		val, valDiags := expr.Value(nil)
		if valDiags.HasErrors() {
			diags = append(diags, valDiags...)
			return nil, nil, diags
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
			if source == formatter.SourceCLI || source == formatter.SourceEnv {
				if !variable.Type.IsPrimitiveType() {
					fakeFilename := fmt.Sprintf("<value for var.%s from source %q>", pbv.Name, source)
					expr, diags = hclsyntax.ParseExpression([]byte(val.AsString()), fakeFilename, hcl.Pos{Line: 1, Column: 1})
					val, valDiags = expr.Value(nil)
					if valDiags.HasErrors() {
						diags = append(diags, valDiags...)
						return nil, nil, diags
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
					Subject: &variable.Range,
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
	for name, variable := range vs {
		v, ok := iv[name]
		if !ok || v == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Unset variable %q", name),
				Detail: "A variable must be set or have a default value; see " +
					"https://www.waypointproject.io/docs/waypoint-hcl/variables/input " +
					"for details.",
				Subject: &variable.Range,
			})
		}
	}
	// Error here if we have them from parsing
	if diags.HasErrors() {
		return nil, nil, diags
	}

	jobVals, diags := getJobValues(vs, iv, salt)
	if diags.HasErrors() {
		return nil, nil, diags
	}
	return iv, jobVals, diags
}

// getJobValues combines the Variable and Value into a VariableRef,
// hashing any 'sensitive' values as SHA256 values with the given salt.
func getJobValues(vs map[string]*Variable, values Values, salt string) (map[string]*pb.Variable_FinalValue, hcl.Diagnostics) {
	varRefs := make(map[string]*pb.Variable_FinalValue, len(values))
	var diags hcl.Diagnostics

	for v, value := range values {
		varRefs[v] = &pb.Variable_FinalValue{}

		// check for sensitive, and salt if so
		if vs[v].Sensitive {
			var sval string
			switch value.Value.Type() {
			case cty.String:
				sval = value.Value.AsString()
			case cty.Bool:
				b := value.Value.True()
				sval = fmt.Sprintf("%t", b)
			case cty.Number:
				var num int64
				err := gocty.FromCtyValue(value.Value, &num)
				// We really shouldn't hit this since we just created the value
				// but for posterity I guess
				if err != nil {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid number",
						Detail:   err.Error(),
					})
					return nil, diags
				}
				sval = fmt.Sprintf("%d", num)
			default:
				// handle any HCL complex types
				bv := hclwrite.TokensForValue(value.Value).Bytes()
				buf := bytes.NewBuffer(bv)
				sval = buf.String()
			}
			// salt shaker
			saltedVal := salt + sval
			h := sha256.Sum256([]byte(saltedVal))
			sval = hex.EncodeToString(h[:])

			varRefs[v].Value = &pb.Variable_FinalValue_Sensitive{Sensitive: sval}
		} else {
			switch value.Value.Type() {
			case cty.String:
				var str string
				err := gocty.FromCtyValue(value.Value, &str)
				if err != nil {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid string",
						Detail:   err.Error(),
					})
					return nil, diags
				}
				varRefs[v].Value = &pb.Variable_FinalValue_Str{Str: str}
			case cty.Bool:
				varRefs[v].Value = &pb.Variable_FinalValue_Bool{Bool: value.Value.True()}
			case cty.Number:
				var num int64
				err := gocty.FromCtyValue(value.Value, &num)
				// We really shouldn't hit this since we just created the value
				// but for posterity I guess
				if err != nil {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid number",
						Detail:   err.Error(),
					})
					return nil, diags
				}
				varRefs[v].Value = &pb.Variable_FinalValue_Num{Num: num}
			default:
				// if it's not a primitive/simple type, we set as bytes here to be later
				// parsed as an hcl expression; any errors at evaluating the hcl type will
				// be handled at that time
				bv := hclwrite.TokensForValue(value.Value).Bytes()
				buf := bytes.NewBuffer(bv)
				varRefs[v].Value = &pb.Variable_FinalValue_Hcl{Hcl: buf.String()}
			}
		}

		source := fromSourceToFV[value.Source]
		varRefs[v].Source = source
	}

	return varRefs, nil
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
			pbv, diags = parseFileValues(f, formatter.SourceVCS)
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
			// if it's not a primitive/simple type, we set as bytes here to be later
			// parsed as an hcl expression; any errors at evaluating the hcl type will
			// be handled at that time
			bv := hclwrite.TokensForValue(val).Bytes()
			buf := bytes.NewBuffer(bv)
			v.Value = &pb.Variable_Hcl{Hcl: buf.String()}
		}

		// Set source
		switch source {
		case formatter.SourceFile:
			v.Source = &pb.Variable_File_{}
		case formatter.SourceVCS:
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
			Subject: &hcl.Range{
				Filename: filename,
			},
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
