package config

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Filename is the default filename for the Waypoint configuration.
const Filename = "waypoint.hcl"

// FindPath looks for our configuration file starting at "start" and
// traversing parent directories until it is found. If it is found, the
// path is returned. If it is not found, an empty string is returned.
// Error will be non-nil only if an error occurred.
//
// If start is empty, start will be the current working directory. If
// filename is empty, it will default to the Filename constant.
func FindPath(start, filename string) (string, error) {
	var err error
	if start == "" {
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	if filename == "" {
		filename = Filename
	}

	for {
		path := filepath.Join(start, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}

		next := filepath.Dir(start)
		if next == start {
			return "", nil
		}

		start = next
	}
}

func (c *Config) LoadPath(path string) error {
	if err := hclsimple.DecodeFile(path, nil, c); err != nil {
		return err
	}

	return c.LoadEnv()
}

func (c *Config) LoadEnv() error {
	return nil
}
