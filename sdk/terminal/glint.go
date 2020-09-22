package terminal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-glint"
	"github.com/olekukonko/tablewriter"
)

type glintUI struct {
	d *glint.Document
}

func GlintUI(ctx context.Context) UI {
	result := &glintUI{
		d: glint.New(),
	}

	go result.d.Render(ctx)

	return result
}

func (ui *glintUI) Close() error {
	return ui.d.Close()
}

func (ui *glintUI) Input(input *Input) (string, error) {
	return "", nil
}

// Interactive implements UI
func (ui *glintUI) Interactive() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}

// Output implements UI
func (ui *glintUI) Output(msg string, raw ...interface{}) {
	msg, style, _ := Interpret(msg, raw...)

	var cs []glint.StyleOption
	switch style {
	case HeaderStyle:
		cs = append(cs, glint.Bold())
		msg = "» " + msg
	case ErrorStyle, ErrorBoldStyle:
		cs = append(cs, glint.Color("red"))
		if style == ErrorBoldStyle {
			cs = append(cs, glint.Bold())
		}

	case WarningStyle:
		cs = append(cs, glint.Color("yellow"))
		if style == WarningBoldStyle {
			cs = append(cs, glint.Bold())
		}

	case SuccessStyle:
		cs = append(cs, glint.Color("green"))
		if style == SuccessBoldStyle {
			cs = append(cs, glint.Bold())
		}

		msg = colorSuccess.Sprint(msg)

	case InfoStyle:
		lines := strings.Split(msg, "\n")
		for i, line := range lines {
			lines[i] = colorInfo.Sprintf("  %s", line)
		}

		msg = strings.Join(lines, "\n")
	}

	ui.d.Append(glint.Finalize(
		glint.Style(
			glint.Text(msg),
			cs...,
		),
	))
}

// NamedValues implements UI
func (ui *glintUI) NamedValues(rows []NamedValue, opts ...Option) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	var buf bytes.Buffer
	tr := tabwriter.NewWriter(&buf, 1, 8, 0, ' ', tabwriter.AlignRight)
	for _, row := range rows {
		switch v := row.Value.(type) {
		case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
			fmt.Fprintf(tr, "  %s: \t%d\n", row.Name, row.Value)
		case float32, float64:
			fmt.Fprintf(tr, "  %s: \t%f\n", row.Name, row.Value)
		case bool:
			fmt.Fprintf(tr, "  %s: \t%v\n", row.Name, row.Value)
		case string:
			if v == "" {
				continue
			}
			fmt.Fprintf(tr, "  %s: \t%s\n", row.Name, row.Value)
		default:
			fmt.Fprintf(tr, "  %s: \t%s\n", row.Name, row.Value)
		}
	}

	tr.Flush()

	ui.d.Append(glint.Finalize(glint.Text(buf.String())))
}

// OutputWriters implements UI
func (ui *glintUI) OutputWriters() (io.Writer, io.Writer, error) {
	return os.Stdout, os.Stderr, nil
}

// Status implements UI
func (ui *glintUI) Status() Status {
	st := newGlintStatus()
	ui.d.Append(st)
	return st
}

func (ui *glintUI) StepGroup() StepGroup {
	ctx, cancel := context.WithCancel(context.Background())
	sg := &glintStepGroup{ctx: ctx, cancel: cancel}
	ui.d.Append(sg)
	return sg
}

// Table implements UI
func (ui *glintUI) Table(tbl *Table, opts ...Option) {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader(tbl.Headers)
	table.SetBorder(false)
	table.SetAutoWrapText(false)

	for _, row := range tbl.Rows {
		colors := make([]tablewriter.Colors, len(row))
		entries := make([]string, len(row))

		for i, ent := range row {
			entries[i] = ent.Value

			color, ok := colorMapping[ent.Color]
			if ok {
				colors[i] = tablewriter.Colors{color}
			}
		}

		table.Rich(entries, colors)
	}

	table.Render()

	ui.d.Append(glint.Finalize(glint.Text(buf.String())))
}
