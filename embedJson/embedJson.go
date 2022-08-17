package embedJson

import (
	"embed"
)

//go:embed gen/*.json
var Files embed.FS
