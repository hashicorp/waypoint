// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// A plugin-agnostic function for end-to-end testing the CLI
// Requires a waypoint server already running, and a project directory
// with an app that can build and deploy.
// Requires WP_BINARY and WP_PROJECT_TEMPLATE_PATH env vars
//
// To run this locally, set up a waypoint server (anywhere), clone the waypoint-examples repo,
// and set WP_PROJECT_TEMPLATE_PATH=/path/to/waypoint-examples/docker/static
func TestCliE2E(t *testing.T) {

	wpBinary = Getenv("WP_BINARY", "waypoint")
	projectTemplatePath := Getenv("WP_PROJECT_TEMPLATE_PATH", "")
	if projectTemplatePath == "" {
		t.Fatalf("Missing required environment variable WP_PROJECT_TEMPLATE_PATH")
	}
	projectFiles, err := ioutil.ReadDir(projectTemplatePath)
	if err != nil {
		t.Fatalf("Failed listing files in project template path %s: %s", projectTemplatePath, err)
	}
	foundWaypointConfig := false
	for _, file := range projectFiles {
		if file.Name() == "waypoint.hcl" {
			foundWaypointConfig = true
			break
		}
	}
	if !foundWaypointConfig {
		t.Fatalf("No waypoint.hcl file in project template path %s", projectTemplatePath)
	}

	require := require.New(t)

	projectName, projectDir, err := SetupTestProject(projectTemplatePath)
	require.NoError(err)
	fmt.Printf("Using project dir %s\n", projectDir)
	fmt.Printf("Using project name %s\n", projectName)

	// Leaving the project dir around for now

	wp := NewBinary(t, wpBinary, projectDir)

	wp.Run("init")

	t.Cleanup(func() {
		// TODO: clean up lingering resources (docker containers, db files, etc)
		//os.RemoveAll(projectDir)
	})

	// Doesn't matter what the output is here, just trying to get to a clean slate

	var deployCount, artifactCount, releaseCount int

	t.Run("Context", func(t *testing.T) {
		table := wp.RunTable("context list")
		require.True(len(table.rows) > 0)

		wp.RunWithOutput("connected successfully", "context verify")

		wp.Run("context inspect")

		// Not testing non-read context commands for now, because we're not creating an
		// isolated context environment.
	})

	t.Run("Up workflow", func(t *testing.T) {
		wp.Run("up")
		artifactCount++
		deployCount++
		wp.RunTableExpectLength(artifactCount, "artifact list")
		wp.RunTableExpectLength(deployCount, "deployment list")
		// Docker plugin has no releaser
	})

	t.Run("Build", func(t *testing.T) {
		wp.Run("build")
		artifactCount++
		wp.RunTableExpectLength(artifactCount, "artifact list")
		wp.RunTableExpectLength(artifactCount, "artifact list-builds")
	})

	t.Run("Deploy", func(t *testing.T) {
		wp.Run("deploy")
		deployCount++
		wp.RunTableExpectLength(deployCount, "deployment list")
		wp.Run("deployment destroy 1")
		deployCount--
		wp.RunTableExpectLength(deployCount, "deployment list")
	})

	t.Run("Release", func(t *testing.T) {
		wp.Run("release")
		// The docker plugin has no releaser
		wp.RunTableExpectLength(releaseCount, "release list")
	})

	t.Run("Fmt", func(t *testing.T) {
		wp.Run("fmt")
	})

	t.Run("Help", func(t *testing.T) {
		_, stderr, _ := wp.RunRaw("help")
		require.True(strings.Contains(stderr, "Usage:"))

		// Try another command too
		_, stderr, _ = wp.RunRaw("up", "-help")
		require.Contains(stderr, "Usage: ")
		require.True(strings.Contains(stderr, "Usage:"))
	})

	t.Run("Config", func(t *testing.T) {
		// Set a config var
		wp.Run("config set insideatest=true")
		wp.RunWithOutput("insideatest", "config get")
	})

	t.Run("Project", func(t *testing.T) {
		// Set a config var
		stdout := wp.Run("project list") // Project list doesn't return a table
		require.True(len(strings.Split(stdout, "\n")) > 0)

		wp.RunWithOutput("Project Name: "+projectName, "project inspect")

		// No easy field to test `project apply` on
	})

	t.Run("Runner profiles", func(t *testing.T) {
		wp.Run("runner profile list") // Has no output right now

		wp.Run("runner profile set -name=e2e-test -plugin-type=docker")

		wp.RunWithOutput("e2e-test", "runner profile list")
	})

	t.Run("Cleanup", func(t *testing.T) {
		wp.Run("destroy -auto-approve")

		wp.RunWithOutput("No deployments found", "deployment list") // Should probably be empty table
		// Artifacts aren't currently destroyed during project destroy
		wp.RunTableExpectLength(0, "release list")
	})
}
