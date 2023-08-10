// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"sort"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type filterOption int

const (
	filterOptionAll filterOption = iota
	filterOptionState
	filterOptionPhyState
	filterOptionOrder
)

type filterFlags struct {
	flagStatusFilter []string
	flagPhysState    string

	order *pb.OperationOrder
}

var stateFiltersMap = map[string]*pb.StatusFilter_Filter{
	"unknown": {
		Filter: &pb.StatusFilter_Filter_State{
			State: pb.Status_UNKNOWN,
		},
	},
	"running": {
		Filter: &pb.StatusFilter_Filter_State{
			State: pb.Status_RUNNING,
		},
	},
	"success": {
		Filter: &pb.StatusFilter_Filter_State{
			State: pb.Status_SUCCESS,
		},
	},
	"error": {
		Filter: &pb.StatusFilter_Filter_State{
			State: pb.Status_ERROR,
		},
	},
}

var knownStates []string

func init() {
	for k := range stateFiltersMap {
		knownStates = append(knownStates, k)
	}

	sort.Strings(knownStates)
}

func stateFlagVar(target *[]string) *flag.EnumVar {
	return &flag.EnumVar{
		Name:   "state",
		Target: target,
		Values: knownStates,
		Usage:  "Filter values to have the given status.",
	}
}

func (ff *filterFlags) statusFilters() []*pb.StatusFilter {
	var st []*pb.StatusFilter

	for _, name := range ff.flagStatusFilter {
		st = append(st, &pb.StatusFilter{
			Filters: []*pb.StatusFilter_Filter{
				stateFiltersMap[name],
			},
		})
	}

	return st
}

var physStateMap = map[string]pb.Operation_PhysicalState{
	"any":       pb.Operation_UNKNOWN,
	"pending":   pb.Operation_PENDING,
	"created":   pb.Operation_CREATED,
	"destroyed": pb.Operation_DESTROYED,
}

func phyStateFlagVar(target *string) *flag.EnumSingleVar {
	var phyStates []string

	for k := range physStateMap {
		phyStates = append(phyStates, k)
	}

	sort.Strings(phyStates)

	return &flag.EnumSingleVar{
		Name:    "physical-state",
		Target:  target,
		Values:  phyStates,
		Default: "created",
		Usage:   "Show values in the given physical states.",
	}
}

func (ff *filterFlags) physState() (pb.Operation_PhysicalState, error) {
	if ff.flagPhysState == "" {
		return pb.Operation_CREATED, nil
	}

	return physStateMap[ff.flagPhysState], nil
}

func (ff *filterFlags) orderOp() *pb.OperationOrder {
	return ff.order
}

func initFilterFlags(set *flag.Sets, ff *filterFlags, opts filterOption) {
	f := set.NewSet("Filter Options")

	if opts == filterOptionAll || opts == filterOptionState {
		f.EnumVar(stateFlagVar(&ff.flagStatusFilter))
	}

	if opts == filterOptionAll || opts == filterOptionPhyState {
		f.EnumSingleVar(phyStateFlagVar(&ff.flagPhysState))
	}

	if opts == filterOptionAll || opts == filterOptionOrder {
		f.EnumSingleVar(&flag.EnumSingleVar{
			Name:   "order-by",
			Target: new(string),
			Usage:  "Order the values by which field.",
			Values: []string{"start-time", "complete-time"},
			SetHook: func(val string) {
				if ff.order == nil {
					ff.order = &pb.OperationOrder{}
				}

				switch val {
				case "start-time":
					ff.order.Order = pb.OperationOrder_START_TIME
				case "complete-time":
					ff.order.Order = pb.OperationOrder_COMPLETE_TIME
				}
			},
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "desc",
			Target: new(bool),
			Usage:  "Sort the values in descending order.",
			SetHook: func(val bool) {
				if !val {
					return
				}

				if ff.order == nil {
					ff.order = &pb.OperationOrder{}
				}

				ff.order.Desc = val
			},
		})

		f.UintVar(&flag.UintVar{
			Name:   "limit",
			Target: new(uint),
			Usage:  "How many values to show.",
			SetHook: func(val uint) {
				if ff.order == nil {
					ff.order = &pb.OperationOrder{}
				}

				ff.order.Limit = uint32(val)
			},
		})
	}
}
