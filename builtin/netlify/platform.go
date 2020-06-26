package netlify

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	netlify "github.com/netlify/open-api/go/porcelain"
	"github.com/skratchdot/open-golang/open"

	"github.com/hashicorp/waypoint/builtin/files"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// Platform is the Platform implementation for Google Cloud Run.
type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// AuthFunc implements component.Authenticator
func (p *Platform) AuthFunc() interface{} {
	return p.Auth
}

// ValidateAuthFunc implements component.Authenticator
func (p *Platform) ValidateAuthFunc() interface{} {
	return p.ValidateAuth
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// Auth retrieves a token and stores it
func (p *Platform) Auth(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	dir *datadir.Component,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Setup API content for netlify, we are not authenticated yet
	clientContext := apiContext("")

	// If the user configured a token, just stop and use that
	if p.config.AccessToken != "" {
		log.Debug("user configured token in access_token config, not authenticating")
		st.Update("Using configured token")
		return nil
	}

	client := netlify.Default

	// Create a ticket to exchange for a secret token
	ticket, err := client.CreateTicket(clientContext, clientID)
	if err != nil {
		return err
	}

	// Authorize in the users browser
	url := fmt.Sprintf("%s/authorize?response_type=ticket&ticket=%s", netlifyUI, ticket.ID)
	if err := open.Start(url); err != nil {
		err = fmt.Errorf("Error opening URL: %s", err)
		return err
	}

	// Blocks until the user proceeds in the browser
	client.WaitUntilTicketAuthorized(clientContext, ticket)
	if err != nil {
		return err
	}

	token, err := client.ExchangeTicket(clientContext, ticket.ID)
	if err != nil {
		return err
	}

	// Persist the token for future runs
	err = persistLocalToken(dir.DataDir(), token.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

// ValidateAuth checks validity of the stored or supplied credential
func (p *Platform) ValidateAuth(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	dir *datadir.Component,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	st.Update("Validating credentials...")

	// If the user configured a token, just stop and use that
	if p.config.AccessToken != "" {
		log.Debug("user configured token in access_token config, not authenticating")
		st.Update("Using configured token")
		return nil
	}

	// Will retrive local token it if exists
	token := retrieveLocalToken(dir.DataDir())

	clientContext := apiContext(token)
	client := netlify.Default

	// Try listing sites to validate auth
	_, err := client.ListSites(clientContext, nil)
	if err != nil {
		return err
	}

	// Auth is valid if we did not error
	return nil
}

// Deploy deploys a set of files to netlify
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	files *files.Files,
	dir *datadir.Component,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	deployment := &Deployment{}
	client := netlify.Default

	// If the user configured a token, just use that, otherwise
	// get the token that should exist because of auth calls
	var token string
	if p.config.AccessToken == "" {
		localToken := retrieveLocalToken(dir.DataDir())
		token = localToken
	}

	// Setup API content for netlify
	clientContext := apiContext(token)

	// We'll update the user in realtime
	st := ui.Status()

	st.Update("Setting up deploy...")

	site := &models.Site{}

	// If the user specifies a site ID, use it
	if p.config.SiteID != "" {
		retrievedSite, err := client.GetSite(clientContext, p.config.SiteID)
		site = retrievedSite
		if err != nil {
			return nil, err
		}
	} else {
		// If the user specified a site name, use that to find or create
		// otherwise, default to the app name
		siteName := src.App
		if p.config.SiteName != "" {
			siteName = p.config.SiteName
		}

		siteSetup := &models.SiteSetup{
			Site: *&models.Site{
				Name: siteName,
			},
		}

		listParams := operations.NewListSitesParams()
		listParams.Name = &siteName
		sites, err := client.ListSites(clientContext, listParams)
		if err != nil {
			return nil, err
		}

		// Create the site if there are no results for that name, otherwise use
		// it.
		if len(sites) == 0 {
			log.Trace("site does not exist, creating site", "site name", siteName)
			st.Update("Creating site")
			createdSite, err := client.CreateSite(clientContext, siteSetup, false)
			site = createdSite
			if err != nil {
				return nil, err
			}
		} else {
			site = sites[0]
			if site.Name != siteName {
				return nil, fmt.Errorf("site returned does not match")
			}
			log.Trace("found site", "site id", site.ID)
		}
	}

	deployment.SiteId = site.ID
	deployOptions := netlify.DeployOptions{
		IsDraft: true,
		Dir:     files.Path,
		SiteID:  site.ID,
	}

	log.Trace("deploying site", "site id", site.ID)
	st.Update("Deploying site")
	deploy, err := client.DeploySite(clientContext, deployOptions)

	if err != nil {
		return nil, fmt.Errorf("error deploying site: %s", err)
	}

	log.Trace("waiting for deploying to become ready", "site id", site.ID)
	st.Update("Waiting for deploy to be ready")
	deploy, err = client.WaitUntilDeployReady(clientContext, deploy, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	// Clear the status
	st.Close()

	deployment.Url = deploy.DeploySslURL
	log.Trace("url available", "url", deploy.DeploySslURL)
	ui.Output("\nURL: %s", deployment.Url, terminal.WithSuccessStyle())

	return deployment, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// SiteID is the site to deploy to
	SiteID string `hcl:"site_id,optional"`
	// SiteName is the name of the site we create. Defaults
	// to the application.
	SiteName string `hcl:"site_name,optional"`
	// AccessToken is the access token to use, will
	// prompt oauth exchange if not specified
	AccessToken string `hcl:"access_token,optional"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
