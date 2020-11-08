package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// AppConfig has the app configuration settings such as env vars.
type AppConfig struct {
	EnvRaw hcl.Expression `hcl:"env,optional"`

	app *App
}

type AppConfigValue struct {
	Key    string
	From   string
	Config map[string]interface{}
}

func (c *AppConfig) Env() (map[string]*AppConfigValue, error) {
	ctx := c.app.ctx

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

		// TODO(mitchellh): for dynamic we'll want to support something else.
		if val.Type() != cty.String {
			var err error
			val, err = convert.Convert(val, cty.String)
			if err != nil {
				return nil, err
			}
		}

		result[key] = &AppConfigValue{
			Key:  key,
			From: "static",
			Config: map[string]interface{}{
				"value": val.AsString(),
			},
		}
	}

	return result, nil
}
