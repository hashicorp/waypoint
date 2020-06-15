package protomappers

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/history"
	pluginhistory "github.com/hashicorp/waypoint/sdk/internal/plugin/history"
	"github.com/hashicorp/waypoint/sdk/internal/pluginargs"
	"github.com/hashicorp/waypoint/sdk/internal/plugincomponent"
	"github.com/hashicorp/waypoint/sdk/internal/pluginterminal"
	pb "github.com/hashicorp/waypoint/sdk/proto"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	Source,
	SourceProto,
	DeploymentConfig,
	DeploymentConfigProto,
	DatadirProject,
	DatadirApp,
	DatadirComponent,
	DatadirProjectProto,
	DatadirAppProto,
	DatadirComponentProto,
	Logger,
	LoggerProto,
	TerminalUI,
	TerminalUIProto,
	HistoryClient,
	HistoryClientProto,
	ReleaseTargets,
	ReleaseTargetsProto,
	LabelSet,
	LabelSetProto,
}

// Source maps Args.Source to component.Source.
func Source(input *pb.Args_Source) (*component.Source, error) {
	var result component.Source
	return &result, mapstructure.Decode(input, &result)
}

// SourceProto
func SourceProto(input *component.Source) (*pb.Args_Source, error) {
	var result pb.Args_Source
	return &result, mapstructure.Decode(input, &result)
}

// DeploymentConfig
func DeploymentConfig(input *pb.Args_DeploymentConfig) (*component.DeploymentConfig, error) {
	var result component.DeploymentConfig
	return &result, mapstructure.Decode(input, &result)
}

func DeploymentConfigProto(input *component.DeploymentConfig) (*pb.Args_DeploymentConfig, error) {
	var result pb.Args_DeploymentConfig
	return &result, mapstructure.Decode(input, &result)
}

// DatadirProject maps *pb.Args_DataDir_Project to *datadir.Project
func DatadirProject(input *pb.Args_DataDir_Project) *datadir.Project {
	dir := datadir.NewBasicDir(input.CacheDir, input.DataDir)
	return &datadir.Project{Dir: dir}
}

func DatadirProjectProto(input *datadir.Project) *pb.Args_DataDir_Project {
	return &pb.Args_DataDir_Project{
		CacheDir: input.CacheDir(),
		DataDir:  input.DataDir(),
	}
}

// DatadirApp maps *pb.Args_DataDir_App to *datadir.App
func DatadirApp(input *pb.Args_DataDir_App) *datadir.App {
	dir := datadir.NewBasicDir(input.CacheDir, input.DataDir)
	return &datadir.App{Dir: dir}
}

func DatadirAppProto(input *datadir.App) *pb.Args_DataDir_App {
	return &pb.Args_DataDir_App{
		CacheDir: input.CacheDir(),
		DataDir:  input.DataDir(),
	}
}

// DatadirComponent maps *pb.Args_DataDir_Component to *datadir.Component
func DatadirComponent(input *pb.Args_DataDir_Component) *datadir.Component {
	dir := datadir.NewBasicDir(input.CacheDir, input.DataDir)
	return &datadir.Component{Dir: dir}
}

func DatadirComponentProto(input *datadir.Component) *pb.Args_DataDir_Component {
	return &pb.Args_DataDir_Component{
		CacheDir: input.CacheDir(),
		DataDir:  input.DataDir(),
	}
}

// Logger maps *pb.Args_Logger to an hclog.Logger
func Logger(input *pb.Args_Logger) hclog.Logger {
	// We use the default logger as the base. Within a plugin we always set
	// it so we can confidently use this. This lets plugins potentially mess
	// with this but that's a risk we have to take.
	return hclog.L().ResetNamed(input.Name)
}

func LoggerProto(log hclog.Logger) *pb.Args_Logger {
	return &pb.Args_Logger{
		Name: log.Name(),
	}
}

// TerminalUI maps *pb.Args_TerminalUI to an hclog.TerminalUI
func TerminalUI(input *pb.Args_TerminalUI) terminal.UI {
	return &pluginterminal.UI{}
}

func TerminalUIProto(ui terminal.UI) *pb.Args_TerminalUI {
	return &pb.Args_TerminalUI{}
}

func ReleaseTargets(input *pb.Args_ReleaseTargets) []component.ReleaseTarget {
	var result []component.ReleaseTarget
	for _, t := range input.Targets {
		result = append(result, component.ReleaseTarget{
			Deployment: &plugincomponent.Deployment{Any: t.Deployment},
			Percent:    uint(t.Percent),
		})
	}

	return result
}

func ReleaseTargetsProto(ts []component.ReleaseTarget) (*pb.Args_ReleaseTargets, error) {
	var result pb.Args_ReleaseTargets
	for _, t := range ts {
		any, err := component.ProtoAny(t.Deployment)
		if err != nil {
			return nil, err
		}

		result.Targets = append(result.Targets, &pb.Args_ReleaseTargets_Target{
			Deployment: any,
			Percent:    uint32(t.Percent),
		})
	}

	return &result, nil
}

func LabelSet(input *pb.Args_LabelSet) *component.LabelSet {
	return &component.LabelSet{
		Labels: input.Labels,
	}
}

func LabelSetProto(labels *component.LabelSet) *pb.Args_LabelSet {
	return &pb.Args_LabelSet{Labels: labels.Labels}
}

// HistoryClient connects to a history.Client served via the plugin interface.
//
// Note these are tested in sdk/internal/plugin via testDynamicFunc.
func HistoryClient(
	ctx context.Context,
	log hclog.Logger,
	input *pb.Args_HistoryClient,
	internal *pluginargs.Internal,
) (history.Client, error) {
	// Create our plugin
	p := &pluginhistory.HistoryPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	conn, err := internal.Broker.Dial(input.StreamId)
	if err != nil {
		return nil, err
	}
	internal.Cleanup.Do(func() { conn.Close() })

	client, err := p.GRPCClient(ctx, internal.Broker, conn)
	if err != nil {
		return nil, err
	}

	return client.(history.Client), nil
}

// HistoryClientProto takes a history.Client and serves it over the plugin interface.
func HistoryClientProto(
	client history.Client,
	log hclog.Logger,
	internal *pluginargs.Internal,
) *pb.Args_HistoryClient {
	// Create our plugin
	p := &pluginhistory.HistoryPlugin{
		Impl:    client,
		Mappers: internal.Mappers,
		Logger:  log,
	}

	id := internal.Broker.NextId()

	// Serve it
	go internal.Broker.AcceptAndServe(id, func(opts []grpc.ServerOption) *grpc.Server {
		server := plugin.DefaultGRPCServer(opts)
		if err := p.GRPCServer(internal.Broker, server); err != nil {
			panic(err)
		}
		return server
	})

	return &pb.Args_HistoryClient{StreamId: id}
}
