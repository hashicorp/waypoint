// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package null

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

// ConfigSourcer implements component.ConfigSourcer
type ConfigSourcer struct {
	config sourceConfig
}

// Config implements component.Configurable
func (cs *ConfigSourcer) Config() (interface{}, error) {
	return &cs.config, nil
}

// ReadFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) ReadFunc() interface{} {
	return cs.read
}

// StopFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) StopFunc() interface{} {
	return cs.stop
}

func (cs *ConfigSourcer) read(
	ctx context.Context,
	log hclog.Logger,
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {
	var results []*pb.ConfigSource_Value
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		results = append(results, result)

		// Decode our configuration
		var singleReq reqConfig
		log.Debug("decoding config", "raw", req.Config)
		if err := mapstructure.WeakDecode(req.Config, &singleReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		// If we have a static value, use it.
		if v := singleReq.StaticValue; v != "" {
			result.Result = &pb.ConfigSource_Value_Value{
				Value: v,
			}
			continue
		}

		// If we have a value from the source config, use that. It is
		// an error if it doesn't exist.
		if key := singleReq.ConfigKey; key != "" {
			value, ok := cs.config.Values[key]
			if !ok {
				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted,
						fmt.Sprintf("server-side config for key %q does not exist", key)).Proto(),
				}
				continue
			}

			result.Result = &pb.ConfigSource_Value_Value{
				Value: value,
			}
			continue
		}

		// Must set one!
		result.Result = &pb.ConfigSource_Value_Error{
			Error: status.New(
				codes.Aborted,
				"One of `static_value` or `config_key` must be set.",
			).Proto(),
		}
	}

	return results, nil
}

func (cs *ConfigSourcer) stop() error {
	// We don't do any background work
	return nil
}

func (cs *ConfigSourcer) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&sourceConfig{}),
		docs.RequestFromStruct(&reqConfig{}),
	)
	if err != nil {
		return nil, err
	}

	doc.Description("Simple configuration values for experimentation or testing.")

	doc.Example(`
config {
  env = {
    "STATIC" = dynamic("null", {
      static_value = "hello"
    })

    "FROM_CONFIG" = dynamic("null", {
      config_key = "foo"
    })
  }
}
`)

	doc.SetRequestField(
		"static_value",
		"A static value to use for the dynamic configuration.",
		docs.Summary(
			"This just returns the value given as the dynamic configuration.",
			"This isn't very \"dynamic\" but it helps to exercise the full dynamic",
			"configuration code paths which can be useful for experimentation or",
			"testing",
			"\n\nThis is not expected to be used in a real-world production system.",
		),
	)

	doc.SetRequestField(
		"config_key",
		"Return a value from the config source configuration.",
		docs.Summary(
			"This looks up the given key in the `values` configuration for the",
			"config sourcer. This can be used to actually test pulling a dynamic",
			"value, except the dynamic value is just Waypoint server-stored.",
			"This is useful for learning about and experimenting with config sourcer",
			"configuration with Waypoint.",
		),
	)

	doc.SetField(
		"values",
		"A mapping of key to value of values that can be pulled with `config_key`.",
		docs.Summary(
			"These values can be sourced using the `config_key` attribute as",
			"as `dynamic` argument. See the `config_key` documentation for",
			"more information on why this is useful.",
		),
	)

	return doc, nil
}

type reqConfig struct {
	StaticValue string `mapstructure:"static_value,optional"`
	ConfigKey   string `mapstructure:"config_key,optional"`
}

type sourceConfig struct {
	Values map[string]string `hcl:"values,optional"`
}
