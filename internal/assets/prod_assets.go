//go:build assetsembedded
// +build assetsembedded

package assets

import (
	"embed"
	"fmt"
	"os"
	"strings"
)

//go:embed ceb/*
var cebFS embed.FS

func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return cebFS.ReadFile(canonicalName)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	file, err := cebFS.Open(canonicalName)
	if err != nil {
		return nil, fmt.Errorf("Asset %s can't read error: %v", name, err)
	}
	return file.Stat()
}
