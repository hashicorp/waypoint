package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func (c *Config) LoadPath(path string) error {
	if err := hclsimple.DecodeFile(path, nil, c); err != nil {
		return err
	}

	return c.LoadEnv()
}

func (c *Config) LoadEnv() error {
	if c.URL != nil {
		if c.URL.Token == "" {
			token := os.Getenv("WAYPOINT_URL_TOKEN")
			if token == "" {
				return fmt.Errorf("URL service configured but no token available (config or env")
			}

			c.URL.Token = token
		}
	}

	return nil
}
