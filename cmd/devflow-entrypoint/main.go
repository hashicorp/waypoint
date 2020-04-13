package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/ceb"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return 1
	}

	hclog.L().SetLevel(hclog.Trace)

	// Run our core logic
	err := ceb.Run(context.Background(),
		ceb.WithEnvDefaults(),
		ceb.WithExec(args),
	)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Error initializing Devflow entrypoint: %s\n", err)
		return 1
	}

	return 0
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		strings.TrimLeftFunc(usageText, unicode.IsSpace),
		os.Args[0])
	flag.PrintDefaults()
}

const usageText = `
Usage: %[1]s [cmd] [args...]

    This the custom entrypoint to support Devflow. It will re-execute any
    command given after configuring the environment for usage with Devflow.

`
