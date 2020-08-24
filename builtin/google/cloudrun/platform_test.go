package cloudrun

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
	c := Config{
		Project: "waypoint-286812",
		Region:  "europe-north1",
		Capacity: &Capacity{
			Memory:                  "512Mi", // max 4Gi
			CPUCount:                2,       // max 2
			MaxRequestsPerContainer: 10,      // default 80, max 80
			RequestTimeout:          500,     // max 900
		},
		AutoScaling: &AutoScaling{
			Max: 10,
		},
		Env: map[string]string{
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

	td, rf := datadir.TestDir(t)
	t.Cleanup(func() {
		rf()
	})

	dir := &datadir.Component{td}
	dc := &component.DeploymentConfig{}
	ui := &StubbedUI{}

	_, err := p.Deploy(
		context.Background(),
		log,
		src,
		img,
		dir,
		dc,
		ui,
	)

	require.NoError(t, err)
}

type StubbedUI struct{}

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
	return &StubbedStatus{}
}

func (sui *StubbedUI) Table(*terminal.Table, ...terminal.Option) {
	return
}

func (sui *StubbedUI) StepGroup() terminal.StepGroup {
	return nil
}

type StubbedStatus struct{}

func (ss *StubbedStatus) Update(msg string)       {}
func (ss *StubbedStatus) Step(status, msg string) {}
func (ss *StubbedStatus) Close() error {
	return nil
}
