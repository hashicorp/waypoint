// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverconfig

import (
	"strconv"
)

const (
	// Default ports. These are strings because we generally are working with
	// strings for the ports since they're part of the address string.
	DefaultGRPCPort = "9701"
	DefaultHTTPPort = "9702"
)

// Client configures a client to connect to a server.
type Client struct {
	Address string `hcl:"address,attr" json:"address"`

	// Tls, if true, will connect to the server with TLS. If TlsSkipVerify
	// is true, the certificate presented by the server will not be validated.
	Tls           bool `hcl:"tls,optional" json:"tls,omitempty"`
	TlsSkipVerify bool `hcl:"tls_skip_verify,optional" json:"tls_skip_verify,omitempty"`

	// AddressInternal is a temporary config to work with local deployments
	// on platforms such as Docker for Mac. We need to discuss a more
	// long term approach to this.
	AddressInternal string `hcl:"address_internal,optional" json:"address_internal,omitempty"`

	// Indicates that we need to present a token to connect to this server.
	RequireAuth bool `hcl:"require_auth,optional" json:"require_path,omitempty"`

	// AuthToken is the token to use to authenticate to the server.
	// Note this will be stored plaintext on disk. You can also use the
	// WAYPOINT_SERVER_TOKEN env var.
	AuthToken string `hcl:"auth_token,optional" json:"auth_token,omitempty"`

	// The platform for where the server is running. Although this option should
	// be required, it's optional to support previously set contexts that did
	// not have a platform.
	Platform string `hcl:"platform,optional" json:"platform,omitempty"`
}

// EnvMap returns a map of environment variables settings
// that will authenticate to the server without a context set.
func (c *Client) EnvMap() map[string]string {
	result := map[string]string{
		"WAYPOINT_SERVER_ADDR":            c.Address,
		"WAYPOINT_SERVER_TLS":             strconv.FormatBool(c.Tls),
		"WAYPOINT_SERVER_TLS_SKIP_VERIFY": strconv.FormatBool(c.TlsSkipVerify),
	}

	if c.RequireAuth {
		result["WAYPOINT_SERVER_TOKEN"] = c.AuthToken
	}

	return result
}

// Env returns a slice of environment variables in key=value settings
// that will authenticate to the server without a context set.
func (c *Client) Env() []string {
	var result []string

	for k, v := range c.EnvMap() {
		result = append(result, k+"="+v)
	}

	return result
}
