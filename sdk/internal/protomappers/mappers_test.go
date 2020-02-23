package protomappers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/pkg/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
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
