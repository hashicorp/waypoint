package protomappers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/internal-shared/mapper"
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

			f, err := mapper.NewFunc(tt.Mapper)
			require.NoError(err)

			raw, err := f.Call(tt.Input...)
			if tt.Error != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Error)
				return
			}
			require.NoError(err)
			require.Equal(tt.Output, raw)
		})
	}
}
