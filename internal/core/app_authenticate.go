package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// CanAuthenticate returns true if the provided component supports authenticating and
// validating authentication for plugins
func (a *App) CanAuthenticate(comp interface{}) bool {
	_, ok := comp.(component.Authenticator)

	return ok
}

// AuthenticateComponent validated authentication for a specific
// component, and if necessary retrieves credentials from the user
func (a *App) AuthenticateComponent(ctx context.Context, comp interface{}) (interface{}, error) {
	auth, ok := comp.(component.Authenticator)
	if !ok {
		return nil, fmt.Errorf("does not implement authenticator")
	}

	validate := func() error {
		_, err := a.callDynamicFunc(ctx,
			a.logger,
			nil,
			auth,
			auth.ValidateAuthFunc(),
		)

		return err
	}

	a.UI.Output("Validating credentials...", terminal.WithHeaderStyle())

	// If validate returns an error, try to auth, otherwise we assume
	// we are valid
	if err := validate(); err != nil {
		a.UI.Output(`There are plugins that require authentication. Waypoint
will guide you through authentication.

`)

		a.UI.Output("Logging in...", terminal.WithHeaderStyle())

		_, err = a.callDynamicFunc(ctx,
			a.logger,
			nil,
			auth,
			auth.AuthFunc(),
		)

		if err != nil {
			return nil, err
		}

		if err := validate(); err != nil {
			return nil, err
		}
	}

	// All is well, continue
	return nil, nil
}
