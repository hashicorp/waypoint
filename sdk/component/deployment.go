package component

// DeploymentConfig is the configuration for the behavior of a deployment.
// Platforms should take this argument and use the value to set the appropriate
// settings for the deployment
type DeploymentConfig struct {
	Id                    string
	ServerAddr            string
	ServerTls             bool
	ServerTlsSkipVerify   bool
	EntrypointInviteToken string
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
		if c.ServerTls {
			results["WAYPOINT_SERVER_TLS"] = "1"
		}
		if c.ServerTlsSkipVerify {
			results["WAYPOINT_SERVER_TLS_SKIP_VERIFY"] = "1"
		}

		// Set our token if we have one
		if c.EntrypointInviteToken != "" {
			results["WAYPOINT_CEB_INVITE_TOKEN"] = c.EntrypointInviteToken
		}
	}

	return results
}
