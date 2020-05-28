package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

const DefaultWaypointRegister = "https://api.alpha.waypoint.run"

type RegisterCommand struct {
	*baseCommand

	registerAddr string

	email   string
	eula    bool
	account bool
	name    string
	token   string
	labels  string

	debug bool
}

type RegisterRequest struct {
	Email      string `json:"email"`
	AcceptEULA bool   `json:"accept_eula"`
	Name       string `json:"name"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type HostnameRequest struct {
	Token    string   `json:"token"`
	Hostname string   `json:"hostname"`
	Labels   []string `json:"labels"`
}

type HostnameResponse struct {
	FQDN  string `json:"fqdn"`
	Error string `json:"error"`
}

func (c *RegisterCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	if c.account {
		if !c.eula {
			c.ui.Output("Pass --accept-eula to confirm you accept the Waypoint URL EULA.", terminal.WithErrorStyle())
			return 1
		}

		var rr RegisterRequest
		rr.Email = c.email
		rr.AcceptEULA = c.eula
		rr.Name = c.name

		var buf bytes.Buffer

		err := json.NewEncoder(&buf).Encode(&rr)
		if err != nil {
			c.ui.Output("Error encoding request to register account: %s", err, terminal.WithErrorStyle())
			return 1
		}

		resp, err := http.Post(DefaultWaypointRegister+"/register", "application/json", &buf)
		if err != nil {
			c.ui.Output("Error requesting account: %s", err, terminal.WithErrorStyle())
			return 1
		}

		var rresp RegisterResponse

		err = json.NewDecoder(resp.Body).Decode(&rresp)
		if err != nil {
			c.ui.Output("Error decoding registration response: %s", err, terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output("Account registered! Use this token to authenticate in future requests:", terminal.WithHeaderStyle())
		c.ui.Output(rresp.Token)

		return 0
	}

	var hr HostnameRequest
	hr.Token = c.token
	hr.Hostname = c.name
	hr.Labels = strings.Split(c.labels, ",")

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(&hr)
	if err != nil {
		c.ui.Output("Error encoding request to register account: %s", err, terminal.WithErrorStyle())
		return 1
	}

	resp, err := http.Post(DefaultWaypointRegister+"/request-hostname", "application/json", &buf)
	if err != nil {
		c.ui.Output("Error requesting hostname: %s", err, terminal.WithErrorStyle())
		return 1
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.Header.Get("Content-Type") == "application/json" {
			var hrr HostnameResponse

			err = json.NewDecoder(resp.Body).Decode(&hrr)
			if err == nil && hrr.Error != "" {
				c.ui.Output("Unknown error requesting hostname: %s", hrr.Error)
				return 1
			}
		}

		c.ui.Output("Unknown error requesting hostname (status code %d)", resp.StatusCode)
		return 1
	}

	var hrr HostnameResponse

	err = json.NewDecoder(resp.Body).Decode(&hrr)
	if err != nil {
		c.ui.Output("Error decoding hostname response: %s", err, terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Successfully requested hostname: %s", hrr.FQDN)

	return 0
}

func (c *RegisterCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		if c.account {
			f.StringVar(&flag.StringVar{
				Name:   "name",
				Target: &c.name,
				Usage:  "Optional name to associate with account.",
			})

			f.StringVar(&flag.StringVar{
				Name:   "email",
				Target: &c.email,
				Usage:  "Email address to associate with account.",
			})

			f.BoolVar(&flag.BoolVar{
				Name:   "accept-eula",
				Target: &c.eula,
				Usage:  "Indicates you accept the usage EULA to use the Waypoint URL Service.",
			})
		} else {
			f.StringVar(&flag.StringVar{
				Name:   "name",
				Target: &c.name,
				Usage:  "The hostname to request.",
			})

			f.StringVar(&flag.StringVar{
				Name:   "token",
				Target: &c.token,
				Usage:  "Token to authenticate with waypoint cluster service (defaults to WAYPOINT_TOKEN env var).",
			})

			f.StringVar(&flag.StringVar{
				Name:    "labels",
				Aliases: []string{"l"},
				Target:  &c.labels,
				Usage:   "Labels to apply to the service.",
			})
		}
	})
}

func (c *RegisterCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RegisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RegisterCommand) Synopsis() string {
	return ""
}

func (c *RegisterCommand) Help() string {
	return ""
}
