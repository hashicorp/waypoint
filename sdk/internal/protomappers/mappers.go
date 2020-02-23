package protomappers

import (
	"github.com/mitchellh/mapstructure"

	"github.com/mitchellh/devflow/sdk/component"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

var All = []interface{}{
	Source,
}

// Source maps Args.Source to component.Source.
func Source(input *pb.Args_Source) (*component.Source, error) {
	var result component.Source
	return &result, mapstructure.Decode(input, &result)
}
