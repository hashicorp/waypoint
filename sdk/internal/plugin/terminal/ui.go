package terminal

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "github.com/hashicorp/waypoint/sdk/proto"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// UIPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the terminal.UI interface.
type UIPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    terminal.UI       // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *UIPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterTerminalUIServiceServer(s, &uiServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
	})
	return nil
}

func (p *UIPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := pb.NewTerminalUIServiceClient(c)
	return terminal.NewCaptureUI(p.Logger, func(lines []*terminal.CaptureLine) error {
		lineArgs := make([]string, len(lines))
		for i, line := range lines {
			lineArgs[i] = line.Line
		}

		_, err := client.Output(context.Background(), &pb.TerminalUI_OutputRequest{
			Lines: lineArgs,
		})
		return err
	}), nil
}

// uiServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type uiServer struct {
	Impl    terminal.UI
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (s *uiServer) Output(
	ctx context.Context,
	req *pb.TerminalUI_OutputRequest,
) (*empty.Empty, error) {
	for _, line := range req.Lines {
		s.Impl.Output(line)
	}

	return &empty.Empty{}, nil
}

var (
	_ plugin.Plugin              = (*UIPlugin)(nil)
	_ plugin.GRPCPlugin          = (*UIPlugin)(nil)
	_ pb.TerminalUIServiceServer = (*uiServer)(nil)
)
