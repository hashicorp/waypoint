package sort

import (
	"reflect"
	"sort"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// ConfigName sorts config variables by name.
type ConfigName []*pb.ConfigVar

func (s ConfigName) Len() int      { return len(s) }
func (s ConfigName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConfigName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ConfigResolution sorts a set of config variables such that if they
// were evaluated in-order to see if they match an environment, you can
// always take the most recent match as the current value. i.e. conflict
// resolution is handled for you in the iteration order.
type ConfigResolution []*pb.ConfigVar

func (s ConfigResolution) Len() int      { return len(s) }
func (s ConfigResolution) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConfigResolution) Less(i, j int) bool {
	// If either is not set, then we can return true or false the ordering
	// doesn't matter versus a nil value.
	if s[i] == nil || s[j] == nil {
		return false
	}
	ti, tj := s[i].Target, s[j].Target

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
	scopeI, scopeJ := reflect.TypeOf(ti.AppScope), reflect.TypeOf(tj.AppScope)
	if scopeI != scopeJ {
		wi, wj := configScopeWeights[scopeI], configScopeWeights[scopeJ]
		return wi < wj
	}

	// 3. If scopes match, the variable that has a label selector is used.
	// The below also solves #4 which sorts based on length of label selector
	// if they both have one.
	return len(s[i].Target.LabelSelector) < len(s[j].Target.LabelSelector)
}

var (
	_ sort.Interface = (ConfigName)(nil)
	_ sort.Interface = (ConfigResolution)(nil)
)

var (
	// configScopeWeights sets the sort weights for comparing config var
	// scopes. The smaller number (i.e. 0 is smaller than 1) is resolved
	// first (meaning is has the lowest precedence, will be chosen last).
	configScopeWeights = map[reflect.Type]int{
		reflect.TypeOf((*pb.ConfigVar_Target_Global)(nil)):      0,
		reflect.TypeOf((*pb.ConfigVar_Target_Project)(nil)):     1,
		reflect.TypeOf((*pb.ConfigVar_Target_Application)(nil)): 2,
	}
)
