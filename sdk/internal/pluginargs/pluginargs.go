// Package pluginargs
package pluginargs

import (
	"github.com/hashicorp/go-plugin"
)

// Broker is the GRPCBroker so that plugins can setup new streams.
type Broker *plugin.GRPCBroker

// Cleanup can be used to register cleanup functions.
type Cleanup struct {
	f func()
}

// Do registers a cleanup function that will be called when the plugin RPC
// call is complete.
func (c *Cleanup) Do(f func()) {
	oldF := c.f
	c.f = func() {
		if oldF != nil {
			defer oldF()
		}
		f()
	}
}

func (c *Cleanup) Close() error {
	if c.f != nil {
		c.f()
	}

	return nil
}
