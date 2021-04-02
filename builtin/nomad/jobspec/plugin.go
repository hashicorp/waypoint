package jobspec

import (
	"github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation for
// the Nomad plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
