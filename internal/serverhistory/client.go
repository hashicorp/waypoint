// Package serverhistory implements the history.Client interface from the SDK
// by calling directly into the backend server.
package serverhistory

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/history"
	"github.com/mitchellh/devflow/sdk/history/convert"
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

	result, err := convert.Component(
		c.MapperSet,
		resp.Deployments,
		cfg.Type,
		(*component.Deployment)(nil),
	)
	if err != nil {
		return nil, err
	}

	return result.([]component.Deployment), nil
}

var _ history.Client = (*Client)(nil)
