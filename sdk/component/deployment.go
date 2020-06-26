package component

import "fmt"

// DeploymentConfig is the configuration for the behavior of a deployment.
// Platforms should take this argument and use the value to set the appropriate
// settings for the deployment
type DeploymentConfig struct {
	Id             string
	ServerAddr     string
	ServerInsecure bool

	// Maps to WAYPOINT_URL_TOKEN in the ceb
	UrlToken string

	// Maps to WAYPOINT_URL_LABELS in the ceb
	UrlLabels string

	// Maps to WAYPOINT_URL_CONTROL_ADDR in the ceb
	UrlControlAddr string
}

// Env returns the environment variables that should be set for the entrypoint
// binary to have the proper configuration.
func (c *DeploymentConfig) Env() map[string]string {
	results := map[string]string{
		"WAYPOINT_DEPLOYMENT_ID": c.Id,
	}

	if c.ServerAddr == "" {
		// If the server is disabled we set this env var. Note that having
		// no address given also causes it to behave the same way.
		results["WAYPOINT_SERVER_DISABLE"] = "1"
	} else {
		// Note the server address.
		results["WAYPOINT_SERVER_ADDR"] = c.ServerAddr
		if c.ServerInsecure {
			results["WAYPOINT_SERVER_INSECURE"] = "1"
		}
	}

	if c.UrlToken != "" {
		results["WAYPOINT_URL_TOKEN"] = c.UrlToken

		labels := c.UrlLabels
		if labels == "" {
			labels = fmt.Sprintf(":deployment=%s", c.Id)
		} else {
			labels = fmt.Sprintf("%s,:deployment=%s", labels, c.Id)
		}

		results["WAYPOINT_URL_LABELS"] = labels

		if c.UrlControlAddr != "" {
			results["WAYPOINT_URL_CONTROL_ADDR"] = c.UrlControlAddr
		}
	}

	return results
}
