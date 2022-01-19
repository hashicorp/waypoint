package main

import (
	"os"

	"github.com/hashicorp/waypoint/internal_nomore/cli"
)

func main() {
	cli.ExposeDocs = true
	os.Exit(cli.Main([]string{"waypoint", "cli-docs"}))
}
