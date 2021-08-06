package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/hashicorp/waypoint/internal/pkg/partial"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// genericConfig represents the `config` stanza that can be placed
// both in the app and at the project level.
type genericConfig struct {
	// internal are variables which can be seen for templating but are not exposed
	// by default to an application or runner.
	InternalRaw hcl.Expression `hcl:"internal,optional"`

	// env are variables that will be exported into the application or runners
	// environment.
	EnvRaw hcl.Expression `hcl:"env,optional"`

	// file are paths that will be written to disk in the context of the application
	// environment.
	FileRaw hcl.Expression `hcl:"file,optional"`

	// Indicates a signal to send the application when config files change.
	FileChangeSignal string `hcl:"file_change_signal,optional"`

	ctx       *hcl.EvalContext    // ctx is the context to use when evaluating
	scopeFunc func(*pb.ConfigVar) // scopeFunc should set the scope for the config var
}

func (c *genericConfig) ConfigVars() ([]*pb.ConfigVar, error) {
	if c == nil {
		return nil, nil
	}

	return c.envVars()
}

var hclEscaper = strings.NewReplacer("${", "$${", "%{", "%%{")

func (c *genericConfig) envVars() ([]*pb.ConfigVar, error) {
	ctx := c.ctx
	ctx = appendContext(ctx, &hcl.EvalContext{
		Functions: map[string]function.Function{
			"configdynamic": configDynamicFunc,
		},
	})
	ctx = finalizeContext(ctx)

	// We're going to build up the variables as we go along using these 4 maps.
	ctx.Variables = map[string]cty.Value{}

	var (
		env      = map[string]cty.Value{}
		internal = map[string]cty.Value{}
		file     = map[string]cty.Value{}
		config   = map[string]cty.Value{}
	)

	// sortVars performs a topological sort of the variables via references, so
	// the pairs can be evaluated top to bottom safely.
	pairs, err := c.sortVars(ctx)
	if err != nil {
		return nil, err
	}

	var result []*pb.ConfigVar
	for _, pair := range pairs {
		key := pair.Name

		// Start building our var
		var newVar pb.ConfigVar
		c.scopeFunc(&newVar)
		newVar.Name = key
		newVar.Internal = pair.Internal
		newVar.NameIsPath = pair.Path

		// Decode the value
		val, diags := pair.Pair.Value.Value(ctx)
		if diags.HasErrors() {
			// Ok, we can't read it's value right now. Let's do a partial evaluation then.
			str, err := partial.EvalExpression(ctx, pair.Pair.Value)
			if err != nil {
				return nil, err
			}

			newVar.Value = &pb.ConfigVar_Static{
				Static: str,
			}

			// We don't advertise these variables in the eval context because
			// we don't want them to be substituted as strings into other variables.
			// If the current variable is referenced by a later variable, we want
			// that to be a normal HCL template expansion of the variable reference,
			// not the contents. Quick example:
			//
			// a = "${g} ${s}"
			// b = "more: ${a}"
			// g = unknown()
			// s = "ok"
			//
			// After running the algorithm, we want b to still be 'more: ${a}', NOT
			// 'more: ${g} ok'. The reason being the 2nd one confuses the escaping
			// as it appears like it might be data that was returned from a file or
			// something.
		} else {
			switch val.Type() {
			case typeDynamicConfig:
				newVar.Value = &pb.ConfigVar_Dynamic{
					Dynamic: val.EncapsulatedValue().(*pb.ConfigVar_DynamicVal),
				}

			default:
				// For non-config val types we try to convert it to a string
				// as a static value.
				var err error
				val, err = convert.Convert(val, cty.String)
				if err != nil {
					return nil, err
				}

				// We have to escape any HCL we find in the string so that we don't
				// evaluate it down-stream.
				// First, we need to check if the value is not null, since we allow
				// `null` defaults for input variables, and a user may forget to
				// provide a value to an input variable
				if val.IsNull() {
					return nil, fmt.Errorf("could not evaluate %q in app config with `null` value", newVar.Name)
				}
				newVar.Value = &pb.ConfigVar_Static{
					Static: hclEscaper.Replace(val.AsString()),
				}

				if pair.Internal {
					internal[pair.Name] = val

					// Because of the nature of the hcl map type, we have to rebuild these
					// each time we modify them.
					config["internal"] = cty.MapVal(internal)
				} else if pair.Path {
					file[pair.Name] = val

					// Because of the nature of the hcl map type, we have to rebuild these
					// each time we modify them.
					config["file"] = cty.MapVal(file)
				} else {
					env[pair.Name] = val

					// Because of the nature of the hcl map type, we have to rebuild these
					// each time we modify them.
					config["env"] = cty.MapVal(env)
				}

				ctx.Variables["config"] = cty.MapVal(config)
			}
		}

		result = append(result, &newVar)
	}

	return result, nil
}

var (
	typeDynamicConfig = cty.Capsule("configval",
		reflect.TypeOf((*pb.ConfigVar_DynamicVal)(nil)).Elem())

	// configDynamicFunc implements the configdynamic() HCL function.
	configDynamicFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "from",
				Type: cty.String,
			},

			{
				Name: "config",
				Type: cty.Map(cty.String),
			},
		},
		Type: function.StaticReturnType(typeDynamicConfig),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var config map[string]string
			if err := gocty.FromCtyValue(args[1], &config); err != nil {
				return cty.NilVal, err
			}

			return cty.CapsuleVal(typeDynamicConfig, &pb.ConfigVar_DynamicVal{
				From:   args[0].AsString(),
				Config: config,
			}), nil
		},
	})
)
