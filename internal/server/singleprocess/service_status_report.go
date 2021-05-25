package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UpsertStatusReport(
	ctx context.Context,
	req *pb.UpsertStatusReportRequest,
) (*pb.UpsertStatusReportResponse, error) {
	result := req.StatusReport

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state.StatusReportPut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertStatusReportResponse{StatusReport: result}, nil
}

func (s *service) ListStatusReports(
	ctx context.Context,
	req *pb.ListStatusReportsRequest,
) (*pb.ListStatusReportsResponse, error) {
	result, err := s.state.StatusReportList(req.Application,
		state.ListWithStatusFilter(req.Status...),
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListStatusReportsResponse{StatusReports: result}, nil
}

func (s *service) GetLatestStatusReport(
	ctx context.Context,
	req *pb.GetLatestStatusReportRequest,
) (*pb.StatusReport, error) {
	r, err := s.state.StatusReportLatest(req.Application, req.Workspace)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// GetStatusReport returns a StatusReport based on ID
func (s *service) GetStatusReport(
	ctx context.Context,
	req *pb.GetStatusReportRequest,
) (*pb.StatusReport, error) {
	r, err := s.state.StatusReportGet(req.Ref)
	if err != nil {
		return nil, err
	}

	return r, nil
}
