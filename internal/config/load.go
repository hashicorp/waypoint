package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func (c *Config) LoadPath(path string) error {
	if err := hclsimple.DecodeFile(path, nil, c); err != nil {
		return err
	}

	return c.LoadEnv()
}

func (c *Config) LoadEnv() error {
	return nil
}
