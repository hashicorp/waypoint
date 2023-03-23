// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"

	"github.com/hashicorp/waypoint/internal/cli"
)

func main() {
	cli.ExposeDocs = true
	os.Exit(cli.Main([]string{"waypoint", "cli-docs"}))
}
