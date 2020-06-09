package netlify

import (
	"context"
	fmt "fmt"

	"github.com/hashicorp/go-hclog"
	netlify "github.com/netlify/open-api/go/porcelain"
	"github.com/skratchdot/open-golang/open"
)

const (
	// Netlify client ID for the Waypoint OAuth 2 app
	clientID  = "c9ae91915154e308fc7d5501fbc1799f27ca314503a25956d93ab790be473636"
	netlifyUI = "https://app.netlify.com"
)

// Authenticate makes API calls and user interactions appropriate to create
// and return an API access token
func Authenticate(
	ctx context.Context,
	log hclog.Logger,
) (string, error) {
	client := netlify.Default

	// Create a ticket to exchange for a secret token
	ticket, err := client.CreateTicket(ctx, clientID)
	if err != nil {
		return "", err
	}

	// Authorize in the users browser
	url := fmt.Sprintf("%s/authorize?response_type=ticket&ticket=%s", netlifyUI, ticket.ID)
	if err := open.Start(url); err != nil {
		err = fmt.Errorf("Error opening URL: %s", err)
		return "", err
	}

	// Blocks until the user proceeds in the browser
	client.WaitUntilTicketAuthorized(ctx, ticket)
	if err != nil {
		return "", err
	}

	token, err := client.ExchangeTicket(ctx, ticket.ID)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
