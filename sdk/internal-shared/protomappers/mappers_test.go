package protomappers

import (
	"testing"

	"github.com/hashicorp/go-argmapper"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/sdk/component"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func TestMappers(t *testing.T) {
	var cases = []struct {
		Name   string
		Mapper interface{}
		Input  []interface{}
		Output interface{}
		Error  string
	}{
		{
			"Source",
			Source,
			[]interface{}{&pb.Args_Source{App: "foo"}},
			&component.Source{App: "foo"},
			"",
		},

		{
			"SourceProto",
			SourceProto,
			[]interface{}{&component.Source{App: "foo"}},
			&pb.Args_Source{App: "foo"},
			"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			f, err := argmapper.NewFunc(tt.Mapper)
			require.NoError(err)

			var args []argmapper.Arg
			for _, input := range tt.Input {
				args = append(args, argmapper.Typed(input))
			}

			result := f.Call(args...)
			if tt.Error != "" {
				require.Error(result.Err())
				require.Contains(result.Err().Error(), tt.Error)
				return
			}
			require.NoError(result.Err())
			require.Equal(tt.Output, result.Out(0))
		})
	}
}
