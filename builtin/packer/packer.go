package packer

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

var Options = []sdk.Option{
	sdk.WithComponents(&ConfigSourcer{}),
}
