package terminal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// BasicUI
type BasicUI struct {
	status *spinnerStatus
}

// Output implements UI
func (ui *BasicUI) Output(msg string, raw ...interface{}) {
	msg, style, w := Interpret(msg, raw...)

	switch style {
	case HeaderStyle:
		msg = colorHeader.Sprintf("==> %s", msg)
	case ErrorStyle:
		msg = colorError.Sprint(msg)
	case WarningStyle:
		msg = colorWarning.Sprint(msg)
	case SuccessStyle:
		msg = colorSuccess.Sprint(msg)
	case InfoStyle:
		lines := strings.Split(msg, "\n")
		for i, line := range lines {
			lines[i] = colorInfo.Sprintf("    %s", line)
		}

		msg = strings.Join(lines, "\n")
	}

	st := ui.status
	if st != nil {
		st.Pause()
		defer st.Start()
	}

	// Write it
	fmt.Fprintln(w, msg)
}

// NamedValues implements UI
func (ui *BasicUI) NamedValues(rows []NamedValue, opts ...Option) {
	cfg := &config{Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	cfg.Writer.Write([]byte{'\n'})

	var buf bytes.Buffer
	tr := tabwriter.NewWriter(&buf, 1, 8, 0, ' ', tabwriter.AlignRight)
	for _, row := range rows {
		fmt.Fprintf(tr, "%s: \t%s\n", row.Name, row.Value)
	}

	tr.Flush()

	colorInfo.Fprintln(cfg.Writer, buf.String())

	cfg.Writer.Write([]byte{'\n'})
}

// OutputWriters implements UI
func (ui *BasicUI) OutputWriters() (io.Writer, io.Writer, error) {
	return os.Stdout, os.Stderr, nil
}

// Status implements UI
func (ui *BasicUI) Status() Status {
	if ui.status == nil {
		ui.status = newSpinnerStatus()
	}

	return ui.status
}

func (ui *BasicUI) StepGroup() StepGroup {
	ctx, cancel := context.WithCancel(context.Background())
	display := NewDisplay(ctx, color.Output)

	return &fancyStepGroup{
		ctx:     ctx,
		cancel:  cancel,
		display: display,
	}
}
