package aci

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/stretchr/testify/require"
)

func TestDeployment(t *testing.T) {
	t.Skip()
	_, err := deploy(t)
	require.NoError(t, err)
}

//func TestRelease(t *testing.T) {
//	t.Skip()
//	rc := ReleaserConfig{}
//
//	r := &Releaser{rc}
//	log := hclog.New(&hclog.LoggerOptions{Level: hclog.Debug})
//	ui := &StubbedUI{log}
//
//	// deploy first
//	d, err := deploy(t)
//	require.NoError(t, err)
//
//	_, err = r.Release(
//		context.Background(),
//		log,
//		ui,
//		d,
//	)
//
//	require.NoError(t, err)
//}

func deploy(t *testing.T) (*Deployment, error) {
	c := Config{
		ResourceGroup: "minecraft",

		Capacity: &Capacity{
			Memory:   512, // max 4Gi
			CPUCount: 2,   // max 2
		},
		StaticEnvVars: map[string]string{
			"foo": "bar",
		},
		Ports: []int{80},

		Volumes: []Volume{
			{
				Name: "vol1",
				Path: "/vol",
				GitRepoVolume: &GitRepoVolume{
					Repository: "github.com/hashicorp/consul",
					Revision:   "v1.8.2",
				},
			},
		},
	}

	p := &Platform{c}

	log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString(os.Getenv("WAYPOINT_LOG_LEVEL"))})

	img := &docker.Image{
		Image: "nginx",
		Tag:   "latest",
	}

	src := &component.Source{
		App: "wpmini",
	}

	td, rf := datadir.TestDir(t)
	t.Cleanup(func() {
		rf()
	})

	dir := &datadir.Component{td}
	dc := &component.DeploymentConfig{}
	ui := &StubbedUI{log}

	return p.Deploy(
		context.Background(),
		log,
		src,
		img,
		dir,
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
