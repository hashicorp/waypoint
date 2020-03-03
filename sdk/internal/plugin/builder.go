package plugin

import (
	"context"
	"fmt"
	"reflect"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/proto"
)

// BuilderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Builder component type.
type BuilderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Builder // Impl is the concrete implementation
	Mappers []*mapper.Func    // Mappers
}

func (p *BuilderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterBuilderServer(s, &builderServer{Impl: p.Impl, Mappers: p.Mappers})
	return nil
}

func (p *BuilderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &builderClient{client: proto.NewBuilderClient(c)}, nil
}

// builderClient is an implementation of component.Builder that
// communicates over gRPC.
type builderClient struct {
	client proto.BuilderClient
}

func (c *builderClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *builderClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *builderClient) BuildFunc() interface{} {
	// Get the build spec
	spec, err := c.client.BuildSpec(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}

	return specToFunc(spec, c.build)
}

func (c *builderClient) build(
	ctx context.Context,
	args dynamicArgs,
) (interface{}, error) {
	// Call our function
	resp, err := c.client.Build(ctx, &proto.Build_Args{Args: args})
	if err != nil {
		return nil, err
	}

	// We return the *any.Any directly.
	return resp.Result, nil
}

// builderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type builderServer struct {
	Impl    component.Builder
	Mappers []*mapper.Func
}

func (s *builderServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *builderServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *builderServer) BuildSpec(
	ctx context.Context,
	args *proto.Empty,
) (*proto.FuncSpec, error) {
	return funcToSpec(s.Impl.BuildFunc(), s.Mappers)
}

func (s *builderServer) Build(
	ctx context.Context,
	args *proto.Build_Args,
) (*proto.Build_Resp, error) {
	// Decode all our arguments
	decoded := make([]interface{}, len(args.Args)+1)
	decoded[0] = ctx
	for idx, arg := range args.Args {
		name, err := ptypes.AnyMessageName(arg)
		if err != nil {
			return nil, err
		}

		typ := protobuf.MessageType(name)
		if typ == nil {
			return nil, fmt.Errorf("cannot decode type: %s", name)
		}

		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		v := reflect.New(typ)
		v.Elem().Set(reflect.Zero(typ))

		if err := ptypes.UnmarshalAny(arg, v.Interface().(protobuf.Message)); err != nil {
			return nil, err
		}

		decoded[idx+1] = v.Interface()
	}

	f, err := mapper.NewFunc(s.Impl.BuildFunc())
	if err != nil {
		return nil, err
	}

	chain, err := f.Chain(s.Mappers, decoded...)
	if err != nil {
		return nil, err
	}

	result, err := chain.Call()
	if err != nil {
		return nil, err
	}

	msg, ok := result.(protobuf.Message)
	if !ok {
		return nil, fmt.Errorf("result of plugin-based function must be a proto.Message, got %T", msg)
	}

	encoded, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, err
	}

	return &proto.Build_Resp{Result: encoded}, nil
}

// appendArgs is a helper to encode a number of protobuf.Message into
// any.Any and add it to the list of dynamicArgs to make it easier to build
// up a dynamic function call.
func appendArgs(args dynamicArgs, ms ...protobuf.Message) (dynamicArgs, error) {
	for _, m := range ms {
		encoded, err := ptypes.MarshalAny(m)
		if err != nil {
			return nil, err
		}

		args = append(args, encoded)
	}

	return args, nil
}

var (
	_ plugin.Plugin                = (*BuilderPlugin)(nil)
	_ plugin.GRPCPlugin            = (*BuilderPlugin)(nil)
	_ proto.BuilderServer          = (*builderServer)(nil)
	_ component.Builder            = (*builderClient)(nil)
	_ component.Configurable       = (*builderClient)(nil)
	_ component.ConfigurableNotify = (*builderClient)(nil)
)
