package config

import (
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// genericConfig represents the `config` stanza that can be placed
// both in the app and at the project level.
type genericConfig struct {
	EnvRaw hcl.Expression `hcl:"env,optional"`

	ctx       *hcl.EvalContext    // ctx is the context to use when evaluating
	scopeFunc func(*pb.ConfigVar) // scopeFunc should set the scope for the config var
}

func (c *genericConfig) ConfigVars() ([]*pb.ConfigVar, error) {
	if c == nil {
		return nil, nil
	}

	return c.envVars()
}

func (c *genericConfig) envVars() ([]*pb.ConfigVar, error) {
	ctx := c.ctx
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

	var result []*pb.ConfigVar
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

		// Start building our var
		var newVar pb.ConfigVar
		c.scopeFunc(&newVar)
		newVar.Name = key

		// Decode the value
		val, diags = pair.Value.Value(ctx)
		if diags.HasErrors() {
			return nil, diags
		}

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

			newVar.Value = &pb.ConfigVar_Static{
				Static: val.AsString(),
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
