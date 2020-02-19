package plugin

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component/mocks"
)

func TestPlugins(t *testing.T) {
	require := require.New(t)

	mock := &mocks.Builder{}
	plugins := Plugins(mock)
	bp := plugins[1]["builder"].(*BuilderPlugin)
	require.Equal(bp.Impl, mock)
}
