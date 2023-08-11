// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *State) AddOnDefinitionPut(ctx context.Context, definition *pb.AddOnDefinition) (*pb.AddOnDefinition, error) {
	return nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnDefinitionUpdate(ctx context.Context, definition *pb.AddOnDefinition, existingDefinition *pb.Ref_AddOnDefinition) (*pb.AddOnDefinition, error) {
	return nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnDefinitionGet(ctx context.Context, definition *pb.Ref_AddOnDefinition) (*pb.AddOnDefinition, error) {
	return nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnDefinitionDelete(ctx context.Context, definition *pb.Ref_AddOnDefinition) error {
	return status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnDefinitionList(ctx context.Context, request *pb.ListAddOnDefinitionsRequest) ([]*pb.AddOnDefinition, *pb.PaginationResponse, error) {
	return nil, nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnPut(ctx context.Context, addOn *pb.AddOn) (*pb.AddOn, error) {
	return nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnUpdate(ctx context.Context, addOn *pb.AddOn) (*pb.AddOn, error) {
	return nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnGet(ctx context.Context, addOn *pb.Ref_AddOn) (*pb.AddOn, error) {
	return nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnDelete(ctx context.Context, addOn *pb.Ref_AddOn) error {
	return status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}

func (s *State) AddOnList(ctx context.Context, request *pb.ListAddOnsRequest) ([]*pb.AddOn, *pb.PaginationResponse, error) {
	return nil, nil, status.Errorf(codes.Unimplemented, "Add On Unimplemented")
}
