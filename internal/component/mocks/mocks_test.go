package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/component"
)

func TestForType(t *testing.T) {
	for typ := range component.TypeMap {
		t.Run(typ.String(), func(t *testing.T) {
			require.NotNil(t, ForType(typ))
		})
	}
}

func TestMock(t *testing.T) {
	require := require.New(t)

	b := &Builder{}
	m := Mock(b)
	require.NotNil(m)
	require.Equal(m, &b.Mock)
}
