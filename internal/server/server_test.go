package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
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
