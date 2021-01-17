package netlify

import (
	"context"
	fmt "fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/go-hclog"
	netlify "github.com/netlify/open-api/go/porcelain"
	netlifyContext "github.com/netlify/open-api/go/porcelain/context"
	"github.com/skratchdot/open-golang/open"
)

const (
	// Netlify client ID for the Waypoint OAuth 2 app
	defaultClientID = "c9ae91915154e308fc7d5501fbc1799f27ca314503a25956d93ab790be473636"
	netlifyUI       = "https://app.netlify.com"
	tokenFilename   = "netlify-access-token"
)

// credentials returns a ClientAuthInfoWriter that
// applies the API token to the authentication header if it
// exists
func credentials(token string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		// todo(pearkes): use a proper user agent
		r.SetHeaderParam("User-Agent", "wp")
		if token != "" {
			r.SetHeaderParam("Authorization", "Bearer "+token)
		}
		return nil
	})
}

// persistLocalToken returns a token from the specified directory
// if the file already exists, this will overwrite it.
func persistLocalToken(dir string, token string) error {
	path := filepath.Join(dir, tokenFilename)
	tokenFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer tokenFile.Close()

	_, err = tokenFile.WriteString(token)

	if err != nil {
		return err
	}

	return nil
}

// retrieveLocalToken returns a token from the specified directory
func retrieveLocalToken(dir string) string {
	path := filepath.Join(dir, tokenFilename)

	// Intentionally not checking the error here as we'll
	// validate the token
	token, _ := ioutil.ReadFile(path)

	return string(token)
}

// apiContext returns context.Context suitable for Netlify
// API operations. If an access token is blank it will return
// an unauthenticated context
func apiContext(accessToken string) context.Context {
	ctx := context.Background()

	ctx = netlifyContext.WithAuthInfo(ctx, credentials(accessToken))

	return ctx
}

// Authenticate makes API calls and user interactions appropriate to create
// and return an API access token
func Authenticate(
	ctx context.Context,
	config Config,
	log hclog.Logger,
) (string, error) {
	client := netlify.Default

	log.Trace("creating netlify ticket")
	clientID := config.ClientID

	if clientID == "" {
		clientID = defaultClientID
	}

	// Create a ticket to exchange for a secret token
	ticket, err := client.CreateTicket(ctx, clientID)
	if err != nil {
		return "", err
	}

	// Authorize in the users browser
	url := fmt.Sprintf("%s/authorize?response_type=ticket&ticket=%s", netlifyUI, ticket.ID)
	if err := open.Start(url); err != nil {
		err = fmt.Errorf("error opening URL: %s", err)
		return "", err
	}

	log.Trace("waiting for authorized user", "ticket id", ticket.ID)
	// Blocks until the user proceeds in the browser
	if _, err := client.WaitUntilTicketAuthorized(ctx, ticket); err != nil {
		return "", err
	}

	log.Trace("exchanging ticket for token", "ticket id", ticket.ID)
	token, err := client.ExchangeTicket(ctx, ticket.ID)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
