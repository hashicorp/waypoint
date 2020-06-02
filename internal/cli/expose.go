package cli

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/horizon/pkg/agent"
	"github.com/hashicorp/horizon/pkg/discovery"
	"github.com/hashicorp/horizon/pkg/pb"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

const DefaultHorizonAddress = "control.alpha.hzn.network"

type ExposeCommand struct {
	*baseCommand

	horizonAddr   string
	horizonToken  string
	horizonLabels string
	horizonPort   string
	tcpService    bool
	debug         bool

	test     bool
	testAddr string
}

func (c *ExposeCommand) authToken() string {
	if c.horizonToken != "" {
		return c.horizonToken
	}

	return os.Getenv("WAYPOINT_TOKEN")
}

func (c *ExposeCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	L := c.Log

	ctx := hclog.WithContext(c.Ctx, L)

	if c.test {
		c.ui.Output("Running test HTTP server: %s", c.testAddr)
		go http.ListenAndServe(":"+c.testAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome to Waypoint! This means you've successfully configured an application to be accessible via the Waypoint URL Service. Next up, your app!\n")
		}))

		<-ctx.Done()
		return 0
	}

	L.Debug("starting agent")

	g, err := agent.NewAgent(L.Named("agent"))
	if err != nil {
		c.ui.Output("Error configuring local interface to waypoint url service: %s", err, terminal.WithErrorStyle())
		return 1
	}

	g.Token = c.authToken()
	if g.Token == "" {
		c.ui.Output("No token available. Use --token or set WAYPOINT_TOKEN.", terminal.WithErrorStyle())
		return 1
	}

	target := c.horizonPort

	labels := pb.ParseLabelSet(c.horizonLabels)

	if !c.tcpService {
		if strings.IndexByte(target, ':') == -1 {
			_, err := strconv.Atoi(target)
			if err == nil {
				target = "127.0.0.1:" + target
			} else {
				target = target + ":80"
			}
		}

		_, err = g.AddService(&agent.Service{
			Type:    "http",
			Labels:  labels,
			Handler: agent.HTTPHandler("http://" + target),
		})

		if err != nil {
			c.ui.Output("Error registering service: %s", err, terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output("registering HTTP service at %s", target)
	} else {
		if strings.IndexByte(target, ':') == -1 {
			_, err := strconv.Atoi(target)
			if err == nil {
				target = "127.0.0.1:" + target
			} else {
				c.ui.Output("Unable to interpret '%s' as TCP target address", target, terminal.WithErrorStyle())
				return 1
			}
		}

		_, err = g.AddService(&agent.Service{
			Type:    "tcp",
			Labels:  labels,
			Handler: agent.TCPHandler(target),
		})

		if err != nil {
			c.ui.Output("Error registering service: %s", err, terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output("registering TCP service at %s", target)
	}

	L.Debug("discovering hubs")

	dc, err := discovery.NewClient(c.horizonAddr)
	if err != nil {
		c.ui.Output("Error connecting to waypoint control service: %s", err, terminal.WithErrorStyle())
		return 1
	}

	L.Debug("refreshing data")

	err = dc.Refresh(ctx)
	if err != nil {
		c.ui.Output("Error discovering network endpoints: %s", err, terminal.WithErrorStyle())
		return 1
	}

	err = g.Start(ctx, dc)
	if err != nil {
		c.ui.Output("Error service traffic: %s", err, terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Serving connections to configured service: %s", c.horizonLabels)

	err = g.Wait(ctx)
	if err != nil {
		c.ui.Output("Error service traffic: %s", err, terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *ExposeCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		if c.test {
			f.StringVar(&flag.StringVar{
				Name:    "port",
				Target:  &c.testAddr,
				Default: "8081",
				Usage:   "Port for the test server to listen on.",
			})

			return
		}

		f.StringVar(&flag.StringVar{
			Name:    "cluster-address",
			Target:  &c.horizonAddr,
			Default: DefaultHorizonAddress,
			Usage:   "Address of the waypoint cluster to expose service on.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "token",
			Target: &c.horizonToken,
			Usage:  "Token to authenticate with waypoint cluster service (defaults to WAYPOINT_TOKEN env var).",
		})

		f.StringVar(&flag.StringVar{
			Name:    "labels",
			Aliases: []string{"l"},
			Target:  &c.horizonLabels,
			Usage:   "Labels to apply to the service.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "address",
			Aliases: []string{"p", "port"},
			Target:  &c.horizonPort,
			Usage:   "Address of service to expose. Either a port number of address:port combo.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "tcp",
			Target: &c.tcpService,
			Usage:  "Indicate that the service is a TCP, rather than HTTP, service.",
		})
	})
}

func (c *ExposeCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ExposeCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ExposeCommand) Synopsis() string {
	return ""
}

func (c *ExposeCommand) Help() string {
	return ""
}
