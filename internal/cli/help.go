package cli

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-glint"
)

// formatHelp takes a raw help string and attempts to colorize it automatically.
func formatHelp(v string) string {
	// Trim the empty space
	v = strings.TrimSpace(v)

	var buf bytes.Buffer
	d := glint.New()
	d.SetRenderer(&glint.TerminalRenderer{
		Output: &buf,

		// We set rows/cols here manually. The important bit is the cols
		// needs to be wide enough so glint doesn't clamp any text and
		// lets the terminal just autowrap it. Rows doesn't make a big
		// difference.
		Rows: 10,
		Cols: 180,
	})

	for _, line := range strings.Split(v, "\n") {
		// Usage: prefix lines
		prefix := "Usage: "
		if strings.HasPrefix(line, prefix) {
			d.Append(glint.Layout(
				glint.Style(
					glint.Text(prefix),
					glint.Color("lightMagenta"),
				),
				glint.Text(line[len(prefix):]),
			).Row())

			continue
		}

		// Alias: prefix lines
		prefix = "Alias: "
		if strings.HasPrefix(line, prefix) {
			d.Append(glint.Layout(
				glint.Style(
					glint.Text(prefix),
					glint.Color("lightMagenta"),
				),
				glint.Text(line[len(prefix):]),
			).Row())

			continue
		}

		// A header line
		if reHelpHeader.MatchString(line) {
			d.Append(glint.Style(
				glint.Text(line),
				glint.Bold(),
			))

			continue
		}

		// Normal line
		d.Append(glint.Text(line))
	}

	d.RenderFrame()
	return buf.String()
}

type helpCommand struct {
	SynopsisText string
	HelpText     string
}

func (c *helpCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (c *helpCommand) Synopsis() string {
	return strings.TrimSpace(c.SynopsisText)
}

func (c *helpCommand) Help() string {
	if c.HelpText == "" {
		return c.SynopsisText
	}

	return c.HelpText
}

func (c *helpCommand) HelpTemplate() string {
	return formatHelp(helpTemplate)
}

var reHelpHeader = regexp.MustCompile(`^[a-zA-Z0-9_-].*:$`)

const helpTemplate = `
Usage: {{.Name}} {{.SubcommandName}} SUBCOMMAND

{{indent 2 (trim .Help)}}{{if gt (len .Subcommands) 0}}

Subcommands:
{{- range $value := .Subcommands }}
    {{ $value.NameAligned }}    {{ $value.Synopsis }}{{ end }}

{{- end }}
`
