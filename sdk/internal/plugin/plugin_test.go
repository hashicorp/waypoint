package plugin

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component/mocks"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/internal-shared/protomappers"
)

func init() {
	// Set our default log level lower for tests
	hclog.L().SetLevel(hclog.Debug)
}

func TestPlugins(t *testing.T) {
	require := require.New(t)

	mock := &mocks.Builder{}
	plugins := Plugins(WithComponents(mock))
	bp := plugins[1]["builder"].(*BuilderPlugin)
	require.Equal(bp.Impl, mock)
}

func testDefaultMappers(t *testing.T) []*mapper.Func {
	var mappers []*mapper.Func
	for _, raw := range protomappers.All {
		f, err := mapper.NewFunc(raw)
		require.NoError(t, err)
		mappers = append(mappers, f)
	}

	return mappers
}
