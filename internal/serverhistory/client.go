// Package serverhistory implements the history.Client interface from the SDK
// by calling directly into the backend server.
package serverhistory

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	servercomponent "github.com/hashicorp/waypoint/internal/server/component"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/history"
)

// Client implements history.Client and provides history using a backend server.
type Client struct {
	APIClient pb.WaypointClient // Client to the API server
	MapperSet []*argmapper.Func // Set of mappers we can use for type conversion
}

// Deployments implements history.Client
func (c *Client) Deployments(ctx context.Context, cfg *history.Lookup) ([]component.Deployment, error) {
	resp, err := c.APIClient.ListDeployments(ctx, &pb.ListDeploymentsRequest{
		Order: &pb.OperationOrder{
			Order: pb.OperationOrder_COMPLETE_TIME,
			Desc:  true,
		},
	})
	if err != nil {
		return nil, err
	}

	result := make([]component.Deployment, 0, len(resp.Deployments))
	for _, v := range resp.Deployments {
		if v.Deployment != nil {
			result = append(result, servercomponent.Deployment(v))
		}
	}

	return result, nil
}

var _ history.Client = (*Client)(nil)
