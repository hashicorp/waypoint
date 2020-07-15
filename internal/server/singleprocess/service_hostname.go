package singleprocess

import (
	"context"
	"strings"

	empty "github.com/golang/protobuf/ptypes/empty"
	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

const (
	hznLabelApp       = "waypoint.hashicorp.com/app"
	hznLabelProject   = "waypoint.hashicorp.com/project"
	hznLabelWorkspace = "waypoint.hashicorp.com/workspace"
)

// TODO: test
func (s *service) CreateHostname(
	ctx context.Context,
	req *pb.CreateHostnameRequest,
) (*pb.CreateHostnameResponse, error) {
	if s.urlClient == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"server doesn't have the URL service enabled")
	}

	// Determine our labels based on our target
	labels := &wphznpb.LabelSet{}
	switch t := req.Target.Target.(type) {
	case *pb.Hostname_Target_Application:
		labels.Labels = append(labels.Labels,
			&wphznpb.Label{Name: hznLabelApp, Value: t.Application.Application.Application},
			&wphznpb.Label{Name: hznLabelProject, Value: t.Application.Application.Project},
			&wphznpb.Label{Name: hznLabelWorkspace, Value: t.Application.Workspace.Workspace},
		)

	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid target type")
	}

	// Build our request
	hostnameReq := &wphznpb.RegisterHostnameRequest{
		// By default we generate a hostname
		Hostname: &wphznpb.RegisterHostnameRequest_Generate{
			Generate: &empty.Empty{},
		},

		Labels: labels,
	}

	// If we have a hostname specified, set it
	if req.Hostname != "" {
		hostnameReq.Hostname = &wphznpb.RegisterHostnameRequest_Exact{
			Exact: req.Hostname,
		}
	}

	// Make the request
	resp, err := s.urlClient.RegisterHostname(ctx, hostnameReq)
	if err != nil {
		return nil, err
	}

	// Extract some data for our result
	hostname := resp.Fqdn
	if idx := strings.Index(hostname, "."); idx != -1 {
		hostname = hostname[:idx]
	}
	labelsMap := map[string]string{}
	for _, label := range labels.Labels {
		labelsMap[label.Name] = label.Value
	}

	return &pb.CreateHostnameResponse{
		Hostname: &pb.Hostname{
			Hostname:     hostname,
			Fqdn:         resp.Fqdn,
			TargetLabels: labelsMap,
		},
	}, nil
}

// TODO: test
func (s *service) ListHostnames(
	ctx context.Context,
	req *pb.ListHostnamesRequest,
) (*pb.ListHostnamesResponse, error) {
	if s.urlClient == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"server doesn't have the URL service enabled")
	}

	resp, err := s.urlClient.ListHostnames(ctx, &wphznpb.ListHostnamesRequest{})
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Hostname, 0, len(resp.Hostnames))
	for _, item := range resp.Hostnames {
		labelsMap := map[string]string{}
		for _, label := range item.Labels.Labels {
			labelsMap[label.Name] = label.Value
		}

		result = append(result, &pb.Hostname{
			Hostname:     item.Hostname,
			Fqdn:         item.Fqdn,
			TargetLabels: labelsMap,
		})
	}

	return &pb.ListHostnamesResponse{Hostnames: result}, nil
}
