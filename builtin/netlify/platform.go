package netlify

import (
	"context"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/go-hclog"
	netlify "github.com/netlify/open-api/go/porcelain"
	netlifyContext "github.com/netlify/open-api/go/porcelain/context"

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

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// netlifyContext returns context.Context suitable for Netlify
// API operations. If an access token is blank it will return
// an unauthenticated context
func (p *Platform) apiContext(accessToken string) context.Context {
	ctx := context.Background()

	apiAuthInfo := func(accessToken string) runtime.ClientAuthInfoWriter {
		return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
			r.SetHeaderParam("User-Agent", "wp-dev")
			if accessToken != "" {
				r.SetHeaderParam("Authorization", "Bearer "+accessToken)
			}
			return nil
		})
	}

	return netlifyContext.WithAuthInfo(ctx, apiAuthInfo(accessToken))
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
	clientContext := p.apiContext("")

	// We'll update the user in realtime
	st := ui.Status()
	defer st.Close()

	// Use configured token, otherwise retrieve one with the user
	token := p.config.AccessToken
	if token == "" {
		st.Update("Logging into your Netlify account")
		token, err := Authenticate(clientContext, log)

		if err != nil {
			return nil, err
		}

		_ = token
	}

	// Setup a new authenticated context
	clientContext = p.apiContext(token)

	st.Update("Setting up deploy")

	// Default siteID to the app name, unless provided
	siteID := src.App
	if p.config.SiteID != "" {
		siteID = p.config.SiteID
	}

	deployOptions := netlify.DeployOptions{
		IsDraft: true,
		Dir:     files.Path,
		SiteID:  siteID,
	}

	log.Trace("deploying site", "site id", siteID)
	st.Update("Deploying site")
	deploy, err := client.DeploySite(ctx, deployOptions)
	if err != nil {
		return nil, err
	}

	log.Trace("waiting for deploying to become ready", "site id", siteID)
	st.Update("Waiting for deploy to be ready")
	deploy, err = client.WaitUntilDeployReady(ctx, deploy, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	deployment.Url = deploy.URL

	return deployment, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// SiteID is the site to deploy to
	SiteID string `hcl:"site_id,optional"`
	// AccessToken is the access token to use, will
	// prompt oauth exchange if not specified
	AccessToken string `hcl:"access_token,optional"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
