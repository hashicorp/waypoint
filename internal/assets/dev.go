//go:generate go-bindata -dev -pkg assets -o dev_assets.go -tags !assetsembedded ceb

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

	for dir != "." {
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
