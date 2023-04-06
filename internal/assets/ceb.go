// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package assets

// CEBArch contains the asset name by architecture. The OS is always
// assumed to be "Linux".
var CEBArch = map[string]string{
	"amd64": "ceb/ceb",
	"arm64": "ceb/ceb-arm64",

	// Docker sometimes uses "aarch64" and sometimes uses "arm64". I don't know
	// why it switches between the two or when but our arm64 build works on
	// both.
	"x86_64":  "ceb/ceb",
	"aarch64": "ceb/ceb-arm64",
}
