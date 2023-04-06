// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package null

import (
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Builder struct {
	config BuilderConfig
}

type BuilderConfig struct{}

type Null struct{}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return func(ui terminal.UI) *emptypb.Empty {
		sg := ui.StepGroup()
		step := sg.Add("performing null build")

		time.Sleep(time.Second * 2)

		step.Update("null build complete")
		step.Done()
		return &emptypb.Empty{}
	}
}

// BuildODRFunc implements component.BuilderODR
func (b *Builder) BuildODRFunc() interface{} {
	return func(ui terminal.UI) *emptypb.Empty {
		sg := ui.StepGroup()
		step := sg.Add("performing null odr build")

		time.Sleep(time.Second * 2)

		step.Update("null odr build complete")
		step.Done()
		return &emptypb.Empty{}
	}
}
