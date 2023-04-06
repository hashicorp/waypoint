// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudrun

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/stretchr/testify/require"
)

// These tests are temporary and just to help speed up the development of the plugin
// will replace when the SDK testing functionality is ready.

func TestDeployment(t *testing.T) {
	t.Skip()
	_, err := deploy(t)
	require.NoError(t, err)
}

func TestRelease(t *testing.T) {
	t.Skip()
	rc := ReleaserConfig{}

	r := &Releaser{rc}
	log := hclog.New(&hclog.LoggerOptions{Level: hclog.Debug})
	ui := &StubbedUI{log}

	// deploy first
	d, err := deploy(t)
	require.NoError(t, err)

	_, err = r.Release(
		context.Background(),
		log,
		ui,
		d,
	)

	require.NoError(t, err)
}

func deploy(t *testing.T) (*Deployment, error) {
	c := Config{
		Project:  "waypoint-286812",
		Location: "europe-north1",
		Capacity: &Capacity{
			Memory:                  512, // max 4Gi
			CPUCount:                2,   // max 2
			MaxRequestsPerContainer: 10,  // default 80, max 80
			RequestTimeout:          500, // max 900
		},
		AutoScaling: &AutoScaling{
			Max: 10,
		},
		StaticEnvVars: map[string]string{
			"foo": "bar",
		},
		Port: 5000,
	}

	p := &Platform{c}

	log := hclog.New(&hclog.LoggerOptions{Level: hclog.Debug})

	img := &docker.Image{
		Image: "gcr.io/waypoint-286812/wpmini",
		Tag:   "latest",
	}

	src := &component.Source{
		App: "wpmini",
	}

	dc := &component.DeploymentConfig{}
	ui := &StubbedUI{log}

	return p.Deploy(
		context.Background(),
		log,
		src,
		img,
		dc,
		ui,
	)
}

type StubbedUI struct{ hclog.Logger }

func (sui *StubbedUI) Input(*terminal.Input) (string, error) {
	return "", nil
}

func (sui *StubbedUI) Interactive() bool {
	return false
}

func (sui *StubbedUI) Output(string, ...interface{}) {
	return
}

func (sui *StubbedUI) NamedValues([]terminal.NamedValue, ...terminal.Option) {
	return
}

func (sui *StubbedUI) OutputWriters() (stdout, stderr io.Writer, err error) {
	return os.Stdout, os.Stderr, nil
}

func (sui *StubbedUI) Status() terminal.Status {
	return &StubbedStatus{sui.Logger}
}

func (sui *StubbedUI) Table(*terminal.Table, ...terminal.Option) {
	return
}

func (sui *StubbedUI) StepGroup() terminal.StepGroup {
	return nil
}

type StubbedStatus struct{ hclog.Logger }

func (ss *StubbedStatus) Update(msg string) {
	ss.Logger.Info(msg)
}
func (ss *StubbedStatus) Step(status, msg string) {
	ss.Logger.Info(msg, "status", status)
}
func (ss *StubbedStatus) Close() error {
	return nil
}
