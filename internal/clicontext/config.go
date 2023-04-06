// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clicontext

import (
	"io"
	"net/url"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

// Config is the structure of the context configuration file. This structure
// can be decoded with hclsimple.DecodeFile.
type Config struct {
	// Server is the configuration to talk to a Waypoint server.
	Server serverconfig.Client `hcl:"server,block"`

	// Workspace represents the current workspace for this context. This value
	// can be set instead of using the -workspace CLI flag. If empty, the
	// default value is "default"
	Workspace string `hcl:"workspace,optional"`
}

// LoadPath loads a context configuration from a filepath.
func LoadPath(path string) (*Config, error) {
	var cfg Config
	err := hclsimple.DecodeFile(path, nil, &cfg)
	return &cfg, err
}

// FromURL parses a URL to a Waypoint server and populates as much of the
// context configuration as possible. This makes a number of assumptions:
//
//   - assumes TLS
//   - assumes TLS skip verify
//
// The skip verify bit is a bad default but it is the most common UX
// getting started and this URL is most commonly used with `waypoint login`
// so we want to provide the smoothest experience there at the expense
// of a slight risk.
func (c *Config) FromURL(v string) error {
	// Ensure our value is a valid URL. This turns example.com into
	// "//example.com" for example. Tests verify this work.
	// See https://github.com/golang/go/issues/19297
	if !strings.Contains(v, "/") {
		v = "//" + v
	}

	u, err := url.Parse(v)
	if err != nil {
		return err
	}

	// Set our defaults
	c.Server.Tls = true
	c.Server.TlsSkipVerify = true
	c.Server.RequireAuth = false

	// Setting the address as the default allows this to work for
	// urls like "foo.com:1234" which url.Parse doesn't handle well at all.
	// We then only override the address if we're sure we got a better value.
	c.Server.Address = v

	// Override
	if u.Host != "" {
		c.Server.Address = u.Host

		// If no port is specified, default to port 9701 which is our default
		// gRPC port for Waypoint installations.
		if idx := strings.IndexByte(u.Host, ':'); idx < 0 {
			c.Server.Address += ":" + serverconfig.DefaultGRPCPort
		}
	}

	// Specifically http will override TLS
	if u.Scheme == "http" {
		c.Server.Tls = false
	}

	return nil
}

// WriteTo implements io.WriterTo and encodes this config as HCL.
func (c *Config) WriteTo(w io.Writer) (int64, error) {
	f := hclwrite.NewFile()
	gohcl.EncodeIntoBody(c, f.Body())
	return f.WriteTo(w)
}
