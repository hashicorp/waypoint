package cli

import (
	"context"
	"log"
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
}

func (c *ExposeCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	err := func(ctx context.Context) error {
		level := hclog.Warn

		if c.debug {
			level = hclog.Trace
		}

		L := hclog.New(&hclog.LoggerOptions{
			Name:  "waypoint",
			Level: level,
		})

		ctx = hclog.WithContext(ctx, L)

		L.Debug("starting agent")

		g, err := agent.NewAgent(L.Named("agent"))
		if err != nil {
			c.ui.Output("Error configuring local interface to waypoint url service: %s", err, terminal.WithErrorStyle())
			return ErrSentinel
		}

		g.Token = c.horizonToken

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
				return ErrSentinel
			}

			c.ui.Output("registering HTTP service at %s", target)
		} else {
			if strings.IndexByte(target, ':') == -1 {
				_, err := strconv.Atoi(target)
				if err == nil {
					target = "127.0.0.1:" + target
				} else {
					c.ui.Output("Unable to interpret '%s' as TCP target address", target, terminal.WithErrorStyle())
					return ErrSentinel
				}
			}

			_, err = g.AddService(&agent.Service{
				Type:    "tcp",
				Labels:  labels,
				Handler: agent.TCPHandler(target),
			})

			if err != nil {
				c.ui.Output("Error registering service: %s", err, terminal.WithErrorStyle())
				return ErrSentinel
			}

			c.ui.Output("registering TCP service at %s", target)
		}

		L.Debug("discovering hubs")

		dc, err := discovery.NewClient(c.horizonAddr)
		if err != nil {
			log.Fatal(err)
		}

		L.Debug("refreshing data")

		err = dc.Refresh(ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = g.Start(ctx, dc)
		if err != nil {
			c.ui.Output("Error service traffic: %s", err, terminal.WithErrorStyle())
			return ErrSentinel
		}

		c.ui.Output("serving connections to configured services")

		err = g.Wait(ctx)
		if err != nil {
			c.ui.Output("Error service traffic: %s", err, terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	}(c.Ctx)

	if err != nil {
		return 1
	}

	return 0
}

func (c *ExposeCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
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
			Target:  &c.horizonLabels,
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
