package ecrpull

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

// Options are the SDK options to use for instantiation.
var Options = []sdk.Option{
	sdk.WithComponents(&Builder{}),
}
