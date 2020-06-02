package assets

import (
	"os"
	"path/filepath"
)

var rootDir string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for len(dir) > 0 {
		path := filepath.Join(dir, "internal/assets")
		if _, err := os.Stat(path); err == nil {
			rootDir = path
			return
		}

		dir = filepath.Dir(dir)
	}

	// Uuuuhhh...
	rootDir = "./internal/assets"
}
