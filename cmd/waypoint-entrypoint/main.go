// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/ceb"
	cebssh "github.com/hashicorp/waypoint/internal/ceb/ssh"
	"github.com/hashicorp/waypoint/internal/pkg/signalcontext"
	ssh "github.com/hashicorp/waypoint/internal/ssh"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Usage = usage
	flag.Parse()

	// This logger always uses debug logging. These logs don't go to
	// the log stream. The entrypoint will create a new logger that respects
	// the WAYPOINT_LOG_LEVEL env var.
	log := hclog.L()
	log.SetLevel(hclog.Debug)

	// Create a context that is cancelled on interrupt
	ctx, closer := signalcontext.WithInterrupt(context.Background(), log)
	defer closer()

	// If WAYPOINT_EXEC_PLUGIN_SSH is set, then the entrypoint is being launched
	// by a platform's exec plugin to provide a place to run an exec command.
	// For instance, the Lambda plugin launches the docker image used by the
	// Lambda function as an ECS task and sets these environment variables. The
	// exec function then uses SSH to connect to the ECS task and runs the users
	// requested exec command inside ECS.
	//
	// Any exec function that wishes to spin up the docker image in a special
	// location is free to use this SSH server functionality to create a context
	// to run exec.

	// port contains which TCP port the ssh server should listen on.
	if port := os.Getenv(ssh.ENVSSHPort); port != "" {

		// Pull the ssh materials out of the environment first
		hostKey, userKey, err := ssh.DecodeFromEnv()
		if err != nil {
			fmt.Fprintf(flag.CommandLine.Output(),
				"Error decoding ssh keys in environment: %s\n", formatError(err))
			return 1
		}

		// The combination of key and hostKey allow both sides to validate that they
		// are communicating with the party they expect to be.

		// Run our core logic
		err = cebssh.RunExecSSHServer(ctx, log, port, hostKey, userKey)
		if err != nil {
			fmt.Fprintf(flag.CommandLine.Output(),
				"Error initializing Waypoint entrypoint: %s\n", formatError(err))
			return 1
		}
		return 0
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return 1
	}

	// Start our debug signal handler
	go debugSignalHandler(ctx, log.Named("debug"))

	// Run our core logic
	err := ceb.Run(ctx,
		ceb.WithEnvDefaults(),
		ceb.WithExec(args))
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Error initializing Waypoint entrypoint: %s\n", formatError(err))
		return 1
	}

	return 0
}

func formatError(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}

	return err.Error()
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		strings.TrimLeftFunc(usageText, unicode.IsSpace),
		os.Args[0])
	flag.PrintDefaults()
}

const usageText = `
Usage: %[1]s [cmd] [args...]

    This the custom entrypoint to support Waypoint. It will re-execute any
    command given after configuring the environment for usage with Waypoint.

`
