// Package serverhistory implements the history.Client interface from the SDK
// by calling directly into the backend server.
package serverhistory

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/ptypes/empty"

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

	raw, err := c.MapperSet.ConvertType(resp.Deployments, cfg.Type)
	if err != nil {
		return nil, err
	}

	rawVal := reflect.ValueOf(raw)
	result := make([]component.Deployment, rawVal.Len())
	for i := 0; i < rawVal.Len(); i++ {
		result[i] = rawVal.Index(i).Interface().(component.Deployment)
	}

	return result, nil
}

var _ history.Client = (*Client)(nil)
