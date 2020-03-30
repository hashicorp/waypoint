// Package serverhistory implements the history.Client interface from the SDK
// by calling directly into the backend server.
package serverhistory

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	servercomponent "github.com/mitchellh/devflow/internal/server/component"
	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/history"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
)

// Client implements history.Client and provides history using a backend server.
type Client struct {
	APIClient pb.DevflowClient // Client to the API server
	MapperSet mapper.Set       // Set of mappers we can use for type conversion
}

// Deployments implements history.Client
func (c *Client) Deployments(ctx context.Context, cfg *history.Lookup) ([]component.Deployment, error) {
	resp, err := c.APIClient.ListDeployments(ctx, &empty.Empty{})
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
