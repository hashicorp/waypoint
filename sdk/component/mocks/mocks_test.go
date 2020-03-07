package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
)

func TestForType(t *testing.T) {
	for typ := range component.TypeMap {
		t.Run(typ.String(), func(t *testing.T) {
			require.NotNil(t, ForType(typ))
		})
	}
}

func TestMock_typed(t *testing.T) {
	require := require.New(t)

	b := &Builder{}
	m := Mock(b)
	require.NotNil(m)
	require.Equal(m, &b.Mock)
}

func TestMock_allTypes(t *testing.T) {
	for typ := range component.TypeMap {
		t.Run(typ.String(), func(t *testing.T) {
			require := require.New(t)
			v := ForType(typ)
			m := Mock(v)
			require.NotNil(m)
		})
	}
}
