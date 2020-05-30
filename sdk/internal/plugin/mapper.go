package plugin

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-argmapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/internal-shared/protomappers"
	funcspec "github.com/hashicorp/waypoint/sdk/internal/funcspec2"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// MapperPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Mapper plugin type.
type MapperPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *MapperPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterMapperServer(s, &mapperServer{
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *MapperPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &MapperClient{
		client: pb.NewMapperClient(c),
		logger: p.Logger,
	}, nil
}

// MapperClient is an implementation of component.Mapper over gRPC.
type MapperClient struct {
	client pb.MapperClient
	logger hclog.Logger
}

// Mappers returns the list of mappers that are supported by this plugin.
func (c *MapperClient) Mappers() ([]*argmapper.Func, error) {
	// Get our list of mapper FuncSpecs
	resp, err := c.client.ListMappers(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, err
	}

	// For each FuncSpec we turn that into a real mapper.Func which calls back
	// into our clien to make an RPC call to generate the proper type.
	var funcs []*argmapper.Func
	for _, spec := range resp.Funcs {
		specCopy := spec

		// We use a closure here to capture spec so that we can provide
		// the correct result type. All we're doing is making our callback
		// call the Map RPC call and return the result/error.
		cb := func(ctx context.Context, args funcspec.Args) (*any.Any, error) {
			resp, err := c.client.Map(ctx, &pb.Map_Request{
				Args:   args,
				Result: specCopy.Result[0].Type,
			})
			if err != nil {
				return nil, err
			}

			return resp.Result, nil
		}

		// Build our funcspec function
		f := funcspec.Func(specCopy, cb, argmapper.Logger(c.logger))

		// Accumulate our functions
		funcs = append(funcs, f)
	}

	return funcs, nil
}

// mapperServer is a gRPC server that implements the Mapper service.
type mapperServer struct {
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (s *mapperServer) ListMappers(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.Map_ListResponse, error) {
	// Go through each mapper and build up our FuncSpecs for each of them.
	var result pb.Map_ListResponse
	for _, m := range s.Mappers {
		fn := m.Func()

		// Skip our built-in protomappers
		if _, ok := protomapperAllMap[reflect.ValueOf(fn).Type()]; ok {
			continue
		}

		spec, err := funcspec.Spec(fn,
			argmapper.ConverterFunc(s.Mappers...),
			argmapper.Logger(s.Logger))
		if err != nil {
			s.Logger.Warn(
				"error converting mapper, will not notify plugin host",
				"func", m.String(),
				"err", err,
			)
			continue
		}

		result.Funcs = append(result.Funcs, spec)
	}

	return &result, nil
}

func (s *mapperServer) Map(
	ctx context.Context,
	args *pb.Map_Request,
) (*pb.Map_Response, error) {
	// Find the output type, which we should know about.
	typ := proto.MessageType(args.Result)
	if typ == nil {
		return nil, status.Newf(
			codes.FailedPrecondition,
			"output type is not known: %s",
			args.Result,
		).Err()
	}

	// Build our function that expects this type as an argument
	// so that we can return it. We do this dynamic function thing so
	// that we can just pretend that this is a function we have so that
	// callDynamicFunc just works.
	f := reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{typ}, []reflect.Type{typ}, false),
		func(args []reflect.Value) []reflect.Value {
			return args
		},
	).Interface()

	// Call it!
	result, err := callDynamicFuncAny2(f, args.Args,
		argmapper.Typed(ctx),
		argmapper.ConverterFunc(s.Mappers...),
	)
	if err != nil {
		return nil, err
	}
	return &pb.Map_Response{Result: result}, nil
}

var (
	_ plugin.Plugin     = (*MapperPlugin)(nil)
	_ plugin.GRPCPlugin = (*MapperPlugin)(nil)
	_ pb.MapperServer   = (*mapperServer)(nil)

	// protomapperAllMap is a set of all the protomapper mappers so
	// that we can easily filter them in ListMappers.
	protomapperAllMap = map[reflect.Type]struct{}{}
)

func init() {
	for _, f := range protomappers.All {
		protomapperAllMap[reflect.TypeOf(f)] = struct{}{}
	}
}
