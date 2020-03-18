package protomappers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/datadir"
	"github.com/mitchellh/devflow/sdk/internal/pluginterminal"
	pb "github.com/mitchellh/devflow/sdk/proto"
	"github.com/mitchellh/devflow/sdk/terminal"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	Source,
	SourceProto,
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
