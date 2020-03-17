package cli

import (
	"strings"

	"github.com/fatih/color"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/terminal"
)

type LogsCommand struct {
	*baseCommand
}

var headerColor = color.New(color.FgCyan)

func (c *LogsCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("exec")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(); err != nil {
		return 1
	}

	cfg := c.cfg
	proj := c.project

	// NOTE(mitchellh): temporary restriction
	if len(cfg.Apps) != 1 {
		proj.UI.Output("only one app is supported at this time", terminal.WithErrorStyle())
		return 1
	}

	// Get our app
	app, err := proj.App(cfg.Apps[0].Name)
	if err != nil {
		c.logError(c.Log, "failed to initialize app", err)
		return 1
	}

	lv, err := app.Logs(ctx)
	if err != nil {
		log.Error("error exec", "error", err)
		return 1
	}

	var pv component.PartitionViewer

	for {
		batch, err := lv.NextLogBatch(ctx)
		if err != nil {
			log.Error("Error retrieving logs", "error", err)
			return 1
		}

		if len(batch) == 0 {
			break
		}

		for _, event := range batch {
			// We use this format rather than regular RFC3339Nano because we use .0
			// instead of .9, which preserves the spacing so the output is always
			// lined up
			ts := event.Timestamp.Format("2006-01-02T15:04:05.000Z07:00")
			short := pv.Short(event.Partition)

			header := headerColor.Sprintf("%s %s: ", ts, short)

			if strings.IndexByte(event.Message, '\n') != -1 {
				parts := strings.Split(event.Message, "\n")

				for _, part := range parts {
					m := header + part
					c.ui.Output(m)
				}
			} else {
				m := header + event.Message
				c.ui.Output(m)
			}
		}
	}

	return 0
}

func (c *LogsCommand) Synopsis() string {
	return ""
}

func (c *LogsCommand) Help() string {
	return ""
}
