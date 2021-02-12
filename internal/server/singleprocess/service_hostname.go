package singleprocess

import (
	"context"
	"reflect"
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
	hznLabelInstance  = "waypoint.hashicorp.com/instance-id"
)

func (s *service) CreateHostname(
	ctx context.Context,
	req *pb.CreateHostnameRequest,
) (*pb.CreateHostnameResponse, error) {
	urlClient := s.urlClient()
	if urlClient == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"server doesn't have the URL service enabled")
	}

	// Determine our labels based on our target
	labels, err := s.hostnameLabelSet(req.Target)
	if err != nil {
		return nil, err
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
	resp, err := urlClient.RegisterHostname(ctx, hostnameReq)
	if err != nil {
		return nil, err
	}

	// Extract some data for our result
	hostname := resp.Fqdn
	if idx := strings.Index(hostname, "."); idx != -1 {
		hostname = hostname[:idx]
	}

	return &pb.CreateHostnameResponse{
		Hostname: &pb.Hostname{
			Hostname:     hostname,
			Fqdn:         resp.Fqdn,
			TargetLabels: s.hostnameLabelSetToMap(labels),
		},
	}, nil
}

func (s *service) ListHostnames(
	ctx context.Context,
	req *pb.ListHostnamesRequest,
) (*pb.ListHostnamesResponse, error) {
	urlClient := s.urlClient()
	if urlClient == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"server doesn't have the URL service enabled")
	}

	// If we have a target given, we get the expected label set for that
	// and build the map.
	var targetMap map[string]string
	if req.Target != nil {
		labels, err := s.hostnameLabelSet(req.Target)
		if err != nil {
			return nil, err
		}

		targetMap = s.hostnameLabelSetToMap(labels)
	}

	resp, err := urlClient.ListHostnames(ctx, &wphznpb.ListHostnamesRequest{})
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Hostname, 0, len(resp.Hostnames))
	for _, item := range resp.Hostnames {
		labelsMap := s.hostnameLabelSetToMap(item.Labels)

		// If we have a target map, then we only include this result if
		// the maps match exactly. In the future we may support subset
		// matching but at this time we do not.
		if targetMap != nil && !reflect.DeepEqual(labelsMap, targetMap) {
			continue
		}

		result = append(result, &pb.Hostname{
			Hostname:     item.Hostname,
			Fqdn:         item.Fqdn,
			TargetLabels: labelsMap,
		})
	}

	return &pb.ListHostnamesResponse{Hostnames: result}, nil
}

func (s *service) DeleteHostname(
	ctx context.Context,
	req *pb.DeleteHostnameRequest,
) (*empty.Empty, error) {
	urlClient := s.urlClient()
	if urlClient == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"server doesn't have the URL service enabled")
	}

	_, err := urlClient.DeleteHostname(ctx, &wphznpb.DeleteHostnameRequest{
		Hostname: req.Hostname,
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// hostnameLabelSet returns the label set for a given target.
func (s *service) hostnameLabelSet(target *pb.Hostname_Target) (*wphznpb.LabelSet, error) {
	labels := &wphznpb.LabelSet{}
	switch t := target.Target.(type) {
	case *pb.Hostname_Target_Application:
		labels.Labels = append(labels.Labels,
			&wphznpb.Label{Name: hznLabelApp, Value: t.Application.Application.Application},
			&wphznpb.Label{Name: hznLabelProject, Value: t.Application.Application.Project},
			&wphznpb.Label{Name: hznLabelWorkspace, Value: t.Application.Workspace.Workspace},
		)

	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid target type")
	}

	return labels, nil
}

// hostnameLabelSetToMap turns a label set into a map.
func (s *service) hostnameLabelSetToMap(labels *wphznpb.LabelSet) map[string]string {
	labelsMap := map[string]string{}
	for _, label := range labels.Labels {
		labelsMap[label.Name] = label.Value
	}

	return labelsMap
}

func (s *service) createHostnameIfNotExist(
	ctx context.Context,
	t *pb.Hostname_Target,
) (*pb.Hostname, error) {
	// First check if we have a matching hostname
	resp, err := s.ListHostnames(ctx, &pb.ListHostnamesRequest{Target: t})
	if err != nil {
		return nil, err
	}

	// If we have any matches, just return the first.
	if len(resp.Hostnames) > 0 {
		return resp.Hostnames[0], nil
	}

	// Create it
	createResp, err := s.CreateHostname(ctx, &pb.CreateHostnameRequest{
		Target: t,
	})
	if err != nil {
		return nil, err
	}

	return createResp.Hostname, nil
}
