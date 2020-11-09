package config

import (
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

// AppConfig has the app configuration settings such as env vars.
type AppConfig struct {
	EnvRaw hcl.Expression `hcl:"env,optional"`

	app *App
}

type AppConfigValue struct {
	Key    string
	From   string
	Config map[string]string
}

func (c *AppConfig) Env() (map[string]*AppConfigValue, error) {
	ctx := c.app.ctx
	ctx = appendContext(ctx, &hcl.EvalContext{
		Functions: map[string]function.Function{
			"configdynamic": configDynamicFunc,
		},
	})
	ctx = finalizeContext(ctx)

	pairs, diags := hcl.ExprMap(c.EnvRaw)
	if diags.HasErrors() {
		return nil, diags
	}

	result := map[string]*AppConfigValue{}
	for _, pair := range pairs {
		// Decode the key. The key must be a string.
		val, diags := pair.Key.Value(ctx)
		if diags.HasErrors() {
			return nil, diags
		}
		if val.Type() != cty.String {
			rng := pair.Key.Range()
			return nil, &hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Summary:     "key must be string",
				Subject:     &rng,
				Expression:  pair.Key,
				EvalContext: ctx,
			}
		}
		key := val.AsString()

		// Decode the value
		val, diags = pair.Value.Value(ctx)
		if diags.HasErrors() {
			return nil, diags
		}

		switch val.Type() {
		case typeDynamicConfig:
			// Good

		default:
			// For non-config val types we try to convert it to a string
			// as a static value.
			var err error
			val, err = convert.Convert(val, cty.String)
			if err != nil {
				return nil, err
			}

			val = cty.CapsuleVal(typeDynamicConfig, &appConfigVal{
				From:   "static",
				Config: map[string]string{"value": val.AsString()},
			})
		}

		configVal := val.EncapsulatedValue().(*appConfigVal)
		result[key] = &AppConfigValue{
			Key:    key,
			From:   configVal.From,
			Config: configVal.Config,
		}
	}

	return result, nil
}

// appConfigVal is the type that config* functions decode into
// as part of a cty.Capsule type when decoding.
type appConfigVal struct {
	From   string
	Config map[string]string
}

var (
	typeDynamicConfig = cty.Capsule("configval",
		reflect.TypeOf((*appConfigVal)(nil)).Elem())

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

			return cty.CapsuleVal(typeDynamicConfig, &appConfigVal{
				From:   args[0].AsString(),
				Config: config,
			}), nil
		},
	})
)
