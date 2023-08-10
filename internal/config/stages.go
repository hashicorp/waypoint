// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"github.com/hashicorp/hcl/v2"
)

type hclStage struct {
	Use    *Use     `hcl:"use,block"`
	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`

	// WorkspaceScoped are workspace-scoped stages.
	WorkspaceScoped []*scopedStage `hcl:"workspace,block"`

	// LabelScoped are label-selector-scoped stages.
	LabelScoped []*scopedStage `hcl:"label,block"`
}

type hclBuild struct {
	Registry *hclStage `hcl:"registry,block"`
	Use      *Use      `hcl:"use,block"`
	Body     hcl.Body  `hcl:",body"`
	Remain   hcl.Body  `hcl:",remain"`

	// WorkspaceScoped are workspace-scoped stages.
	WorkspaceScoped []*scopedStage `hcl:"workspace,block"`

	// LabelScoped are label-selector-scoped stages.
	LabelScoped []*scopedStage `hcl:"label,block"`
}

// scopedStage is used within hclStage for workspace/label scoping.
type scopedStage struct {
	// Scope is the label for the block. This is reused for both workspace
	// and label scoped variables so this could be either of those.
	Scope string `hcl:",label"`

	// Same as hclStage
	Use    *Use     `hcl:"use,block"`
	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
}

// Build are the build settings.
type Build struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`

	// This should not be used directly. This is here for validation.
	// Instead, use App.Registry().
	Registry *Registry `hcl:"registry,block"`

	// Unused for practical reasons, but we need this here so that
	// the decoding validates successfully (HCL doesn't error of
	// unexpected things).
	WorkspaceScoped []*scopedStage `hcl:"workspace,block"`
	LabelScoped     []*scopedStage `hcl:"label,block"`

	ctx *hcl.EvalContext
}

// Registry are the registry settings.
type Registry struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`

	// Unused, see Build
	WorkspaceScoped []*scopedStage `hcl:"workspace,block"`
	LabelScoped     []*scopedStage `hcl:"label,block"`

	ctx *hcl.EvalContext
}

// Deploy are the deploy settings.
type Deploy struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`

	// Unused, see Build
	WorkspaceScoped []*scopedStage `hcl:"workspace,block"`
	LabelScoped     []*scopedStage `hcl:"label,block"`

	ctx *hcl.EvalContext
}

// Release are the release settings.
type Release struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`

	// Unused, see Build
	WorkspaceScoped []*scopedStage `hcl:"workspace,block"`
	LabelScoped     []*scopedStage `hcl:"label,block"`

	ctx *hcl.EvalContext
}

// Hook is the configuration for a hook that runs at specified times.
type Hook struct {
	When      string   `hcl:"when,attr"`
	Command   []string `hcl:"command,attr"`
	OnFailure string   `hcl:"on_failure,optional"`
}

func (h *Hook) ContinueOnFailure() bool {
	return h.OnFailure == "continue"
}

func (b *Step) hclContext() *hcl.EvalContext     { return b.ctx }
func (b *Build) hclContext() *hcl.EvalContext    { return b.ctx }
func (b *Registry) hclContext() *hcl.EvalContext { return b.ctx }
func (b *Deploy) hclContext() *hcl.EvalContext   { return b.ctx }
func (b *Release) hclContext() *hcl.EvalContext  { return b.ctx }
