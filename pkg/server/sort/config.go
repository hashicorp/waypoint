// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sort

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/waypoint/pkg/config/funcs"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/mitchellh/pointerstructure"
	"github.com/zclconf/go-cty/cty"
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

// configRunnerSet splits a set of config vars into a merge set depending
// on priority to match a runner.
func configRunnerSet(
	set []*pb.ConfigVar,
	req *pb.Ref_RunnerId,
) ([][]*pb.ConfigVar, error) {
	// Results go into two buckets
	result := make([][]*pb.ConfigVar, 2)
	const (
		idxAny = 0
		idxId  = 1
	)

	// Go through the iterator and accumulate the results
	for _, current := range set {
		if current.Target.Runner == nil {
			// We are not a config for a runner.
			continue
		}

		idx := -1
		switch ref := current.Target.Runner.Target.(type) {
		case *pb.Ref_Runner_Any:
			idx = idxAny

		case *pb.Ref_Runner_Id:
			idx = idxId

			// We need to match this ID
			if ref.Id.Id != req.Id {
				continue
			}

		default:
			return nil, fmt.Errorf("config has unknown target type: %T", current.Target.Runner.Target)
		}

		result[idx] = append(result[idx], current)
	}

	return result, nil
}

// OrderRequest is the input to OrderVariables, with all the information
// to properly calculate the order to evalutate each variable.
type OrderRequest struct {
	Set       [][]*pb.ConfigVar
	Merge     bool
	Runner    *pb.Ref_RunnerId
	Workspace *pb.Ref_Workspace
	Labels    map[string]string
}

// OrderVariables considers the data in the OrderRequest and calculates
// the correct order to evalutate each ConfigVar, and then returns
// the vars in that order. It also filters variables based on the values
// in the OrderRequest, depending on what is set.
func OrderVariables(req *OrderRequest) ([]*pb.ConfigVar, error) {
	mergeSet := req.Set

	// Sort all of our merge sets by the resolution rules
	for _, set := range mergeSet {
		sort.Sort(ConfigResolution(set))
	}

	// If we have a runner set, then we want to filter all our config vars
	// by runner. This is more complex than that though, because tighter
	// scoped runner refs should overwrite weaker scoped (i.e. ID-ref overwrites
	// Any-ref). So we have to split our merge set from <X, Y> to
	// <X_any, X_id, Y_any, Y_id> so it merges properly later.
	if req.Runner != nil {
		var newMergeSet [][]*pb.ConfigVar
		for _, set := range mergeSet {
			splitSets, err := configRunnerSet(set, req.Runner)
			if err != nil {
				return nil, err
			}

			newMergeSet = append(newMergeSet, splitSets...)
		}

		mergeSet = newMergeSet
	} else {
		// If runner isn't set, then we want to ensure we're not getting
		// any runner env vars.
		for _, set := range mergeSet {
			for i, v := range set {
				if v == nil {
					continue
				}

				if v.Target.Runner != nil {
					set[i] = nil
				}
			}
		}
	}

	// Filter based on the workspace if we have it set.
	if req.Workspace != nil {
		for _, set := range mergeSet {
			for i, v := range set {
				if v == nil {
					continue
				}

				if v.Target.Workspace != nil &&
					!strings.EqualFold(v.Target.Workspace.Workspace, req.Workspace.Workspace) {
					set[i] = nil
				}
			}
		}
	}

	// Filter by labels
	ctyMap := cty.MapValEmpty(cty.String)
	if len(req.Labels) > 0 {
		mapValues := map[string]cty.Value{}
		for k, v := range req.Labels {
			mapValues[k] = cty.StringVal(v)
		}
		ctyMap = cty.MapVal(mapValues)
	}

	for _, set := range mergeSet {
		for i, v := range set {
			if v == nil {
				continue
			}

			// If there is no selector, ignore.
			if v.Target.LabelSelector == "" {
				continue
			}

			// Use our selectormatch HCL function for equal logic
			result, err := funcs.SelectorMatch(ctyMap, cty.StringVal(v.Target.LabelSelector))
			if errors.Is(err, pointerstructure.ErrNotFound) {
				// this means that the label selector contains a label
				// that isn't set, this means we do not match.
				err = nil
				result = cty.BoolVal(false)
			}
			if err != nil {
				return nil, err
			}

			if result.False() {
				set[i] = nil
			}
		}
	}

	// If we aren't merging, then we're done. We just flatten the list.
	if !req.Merge {
		var result []*pb.ConfigVar
		for _, set := range mergeSet {
			for _, v := range set {
				if v != nil {
					result = append(result, v)
				}
			}
		}
		sort.Sort(ConfigName(result))
		return result, nil
	}

	// Merge our merge set
	merged := make(map[string]*pb.ConfigVar)
	for _, set := range mergeSet {
		for _, v := range set {
			// Ignore nil since those are filtered out values.
			if v == nil {
				continue
			}

			merged[v.Name] = v
		}
	}

	result := make([]*pb.ConfigVar, 0, len(merged))
	for _, v := range merged {
		result = append(result, v)
	}

	sort.Sort(ConfigName(result))

	return result, nil

}
