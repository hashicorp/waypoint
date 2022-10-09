package data

import (
	"embed"
	"strings"
)

//go:embed init.tpl.hcl
var dataFS embed.FS

func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return dataFS.ReadFile(canonicalName)
}
