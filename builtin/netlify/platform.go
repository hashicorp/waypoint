package netlify

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	netlify "github.com/netlify/open-api/go/porcelain"

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

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Default siteID to the app name, unless provided
	siteID := src.App
	if p.config.SiteID != "" {
		siteID = p.config.SiteID
	}

	deployOptions := netlify.DeployOptions{
		IsDraft: true,
		Dir:     files.GetDirectory(),
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
	SiteID string `hcl:"site_id,attr"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
