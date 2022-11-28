package sort

import (
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"reflect"
	"sort"
)

type ConfigSource []*pb.ConfigSource

var (
	_ sort.Interface = (ConfigSource)(nil)
)

func (s ConfigSource) Len() int      { return len(s) }
func (s ConfigSource) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConfigSource) Less(i, j int) bool {
	// If either is not set, then we can return true or false the ordering
	// doesn't matter versus a nil value.
	if s[i] == nil || s[j] == nil {
		return false
	}
	ti, tj := s[i], s[j]

	// 1. If a workspace is set on one but not the other, the variable
	//    with the workspace sorts higher than no workspace.
	var wsi, wsj string
	if v := ti.Workspace; v != nil {
		wsi = v.Workspace
	}
	if v := tj.Workspace; v != nil {
		wsj = v.Workspace
	}

	// This boolean statement can be confusing, but in english:
	// if the workspace targets set are not equal and at least one is empty.
	// We know both aren't empty (because they're not equal).
	// This means that resolution rule #1 takes effect.
	if wsi != wsj && (wsi == "" || wsj == "") {
		return wsi == ""
	}

	// 2. The most specific "scope" is used: app over project over global.
	// This only applies so long as the scope types do not match.
	scopeI, scopeJ := reflect.TypeOf(ti.Scope), reflect.TypeOf(tj.Scope)
	if scopeI != scopeJ {
		wi, wj := configSourceScopeWeights[scopeI], configSourceScopeWeights[scopeJ]
		return wi < wj
	}

	// TODO: Support label scoping
	return false
}

var (
	// configScopeWeights sets the sort weights for comparing config var
	// scopes. The smaller number (i.e. 0 is smaller than 1) is resolved
	// first (meaning is has the lowest precedence, will be chosen last).
	configSourceScopeWeights = map[reflect.Type]int{
		reflect.TypeOf((*pb.ConfigSource_Global)(nil)):      0,
		reflect.TypeOf((*pb.ConfigSource_Project)(nil)):     1,
		reflect.TypeOf((*pb.ConfigSource_Application)(nil)): 2,
	}
)
