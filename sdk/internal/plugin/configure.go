package plugin

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/protostructure"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/sdk/component"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// configStruct is the shared helper to implement the ConfigStruct RPC call
// for components. The logic is the same regardless of component so this can
// be called instead.
func configStruct(impl interface{}) (*pb.Config_StructResp, error) {
	c, ok := impl.(component.Configurable)

	// If Configurable isn't implemented, we just return an empty response.
	// The nil struct signals to the receiving side that this component
	// is not configurable.
	if !ok {
		return &pb.Config_StructResp{}, nil
	}

	v, err := c.Config()
	if err != nil {
		return nil, err
	}

	s, err := protostructure.Encode(v)
	if err != nil {
		return nil, err
	}

	return &pb.Config_StructResp{Struct: s}, nil
}

// configStructCall is the shared helper to call the ConfigStruct RPC call
// and return the proper struct value for decoding configuration.
func configStructCall(ctx context.Context, c configurableClient) (interface{}, error) {
	resp, err := c.ConfigStruct(ctx, &empty.Empty{})

	// If we had a failure receiving the configuration struct, then
	// panic because this should never happen. In the future maybe we can
	// support an error return value.
	if err != nil {
		return nil, err
	}

	// If we have no struct, then we have no value so return nil
	if resp.Struct == nil {
		return nil, nil
	}

	result, err := protostructure.New(resp.Struct)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// configure is the shared helper to implement the Configure RPC call.
func configure(impl interface{}, req *pb.Config_ConfigureRequest) (*empty.Empty, error) {
	c, ok := impl.(component.Configurable)

	// This should never happen but if it does just do nothing. This
	// should never happen because prior to this ever being called, our core
	// calls ConfigStruct and if we return nil then we don't configure anything.
	if !ok {
		return &empty.Empty{}, nil
	}

	// Get our value that we can decode into
	v, err := c.Config()
	if err != nil {
		return nil, err
	}

	// Decode our JSON value directly into our structure.
	if err := json.Unmarshal(req.Json, v); err != nil {
		return nil, err
	}

	// If our client also implements the notify interface, call that.
	if cn, ok := c.(component.ConfigurableNotify); ok {
		if err := cn.ConfigSet(v); err != nil {
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

// configureCall calls the Configure RPC endpoint.
func configureCall(ctx context.Context, c configurableClient, v interface{}) error {
	jsonv, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = c.Configure(ctx, &pb.Config_ConfigureRequest{
		Json: jsonv,
	})
	return err
}

// configurableClient is the interface implemented by all gRPC services that
// have the configuration RPC methods. We use this with the helpers above
// to extract shared logic for component configuration.
type configurableClient interface {
	ConfigStruct(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.Config_StructResp, error)
	Configure(context.Context, *pb.Config_ConfigureRequest, ...grpc.CallOption) (*empty.Empty, error)
}
