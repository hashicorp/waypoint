// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package null

import (
	"time"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/builtin/k8s"
)

type Releaser struct {
	config ReleaserConfig
}

type ReleaserConfig struct{}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	// TODO: we should implement a null release proto with the URL function at some point,
	// but a k8s release type is OK for now.
	return func(ui terminal.UI) *k8s.Release {
		sg := ui.StepGroup()
		step := sg.Add("performing null release")

		time.Sleep(time.Second * 2)

		step.Update("null release complete")
		step.Done()
		return &k8s.Release{}
	}
}

// DestroyFunc implements component.Destroyer
func (r *Releaser) DestroyFunc() interface{} {
	return func(ui terminal.UI) error {
		sg := ui.StepGroup()
		step := sg.Add("performing null release destroy")

		time.Sleep(time.Second * 1)

		step.Update("null release destroy complete")
		step.Done()
		return nil
	}
}

// StatusFunc implements component.Status
func (r *Releaser) StatusFunc() interface{} {
	return func(ui terminal.UI) *sdk.StatusReport {
		sg := ui.StepGroup()
		step := sg.Add("performing null release status")

		time.Sleep(time.Second * 1)

		step.Update("null release status")
		step.Done()
		return &sdk.StatusReport{}
	}
}

//func (e *TODO_RELEASER_PROTO_MESSAGE_STRUCT) URL() string { return "empty release url" }
