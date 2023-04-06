// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	goplugin "github.com/hashicorp/go-plugin"
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

// ParseReattachPlugins parses information on reattaching to plugins out of a JSON-encoded environment variable.
// Example input: WP_REATTACH_PLUGINS='{"pack":{"Protocol":"grpc","ProtocolVersion":1,"Pid":24025,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/ns/grk8kk196_106v37w9hk8hxm0000gq/T/plugin047564716"}}}'
func ParseReattachPlugins(in string) (map[string]*goplugin.ReattachConfig, error) {
	reattachConfigs := map[string]*goplugin.ReattachConfig{}
	if in != "" {
		in = strings.TrimRight(in, "'")
		in = strings.TrimLeft(in, "'")
		var m map[string]sdk.ReattachConfig
		err := json.Unmarshal([]byte(in), &m)
		if err != nil {
			return reattachConfigs, fmt.Errorf("Invalid format for WP_REATTACH_PROVIDERS: %w", err)
		}
		for p, c := range m {
			var addr net.Addr
			switch c.Addr.Network {
			case "unix":
				addr, err = net.ResolveUnixAddr("unix", c.Addr.String)
				if err != nil {
					return reattachConfigs, fmt.Errorf("Invalid unix socket path %q for %q: %w", c.Addr.String, p, err)
				}
			case "tcp":
				addr, err = net.ResolveTCPAddr("tcp", c.Addr.String)
				if err != nil {
					return reattachConfigs, fmt.Errorf("Invalid TCP address %q for %q: %w", c.Addr.String, p, err)
				}
			default:
				return reattachConfigs, fmt.Errorf("Unknown address type %q for %q", c.Addr.String, p)
			}
			reattachConfigs[p] = &goplugin.ReattachConfig{
				Protocol:        goplugin.Protocol(c.Protocol),
				ProtocolVersion: c.ProtocolVersion,
				Pid:             c.Pid,
				Test:            c.Test,
				Addr:            addr,
			}
		}
	}
	return reattachConfigs, nil
}
