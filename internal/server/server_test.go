package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/mitchellh/devflow/internal/server/gen"
	"github.com/mitchellh/devflow/sdk/component"
)

func TestComponentEnum(t *testing.T) {
	for idx, name := range pb.Component_Type_name {
		// skip the invalid value
		if idx == 0 {
			continue
		}

		typ := component.Type(idx)
		require.Equal(t, strings.ToUpper(typ.String()), strings.ToUpper(name))
	}
}
