// Package cloudrun contains components for deploying to Google Cloud Run.
package cloudrun

import (
	"github.com/mitchellh/devflow/sdk"
)

//go:generate go-bindata -fs -nomemcopy -nometadata -pkg cloudrun -prefix data/ data/...

// Options are the SDK options to use for instantiation for
// the Google Cloud Run plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
}
