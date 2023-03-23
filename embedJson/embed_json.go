// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package embedJson

import (
	"embed"
)

// This file embeds all of the files in /gen/ into the Files variable below. The comment above
// tells the compiler where to find these files. Once waypoint is built, these files are avaliable
// through this package at runtime.

//go:embed gen/*.json
var Files embed.FS
