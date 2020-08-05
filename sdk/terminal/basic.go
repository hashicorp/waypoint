package terminal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bgentry/speakeasy"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// basicUI
type basicUI struct {
	ctx    context.Context
	status *spinnerStatus
}

// Returns a UI which will write to the current processes
// stdout/stderr.
func ConsoleUI(ctx context.Context) UI {
	return &basicUI{ctx: ctx}
}

// Input implements UI
func (ui *basicUI) Input(input *Input) (string, error) {
	var buf bytes.Buffer

	// Write the prompt, add a space.
	ui.Output(input.Prompt, WithStyle(input.Style), WithWriter(&buf))
	fmt.Fprint(color.Output, strings.TrimRight(buf.String(), "\r\n"))
	fmt.Fprint(color.Output, " ")

	// Ask for input in a go-routine so that we can ignore it.
	errCh := make(chan error, 1)
	lineCh := make(chan string, 1)
	go func() {
		var line string
		var err error
		if input.Secret && isatty.IsTerminal(os.Stdin.Fd()) {
			line, err = speakeasy.Ask("")
		} else {
			r := bufio.NewReader(os.Stdin)
			line, err = r.ReadString('\n')
		}
		if err != nil {
			errCh <- err
			return
		}

		lineCh <- strings.TrimRight(line, "\r\n")
	}()

	select {
	case err := <-errCh:
		return "", err
	case line := <-lineCh:
		return line, nil
	case <-ui.ctx.Done():
		// Print newline so that any further output starts properly
		fmt.Fprintln(color.Output)
		return "", ui.ctx.Err()
	}
}

// Interactive implements UI
func (ui *basicUI) Interactive() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}

// Output implements UI
func (ui *basicUI) Output(msg string, raw ...interface{}) {
	msg, style, w := Interpret(msg, raw...)

	switch style {
	case HeaderStyle:
		msg = colorHeader.Sprintf("\n==> %s", msg)
	case ErrorStyle:
		msg = colorError.Sprint(msg)
	case ErrorBoldStyle:
		msg = colorErrorBold.Sprint(msg)
	case WarningStyle:
		msg = colorWarning.Sprint(msg)
	case WarningBoldStyle:
		msg = colorWarningBold.Sprint(msg)
	case SuccessStyle:
		msg = colorSuccess.Sprint(msg)
	case SuccessBoldStyle:
		msg = colorSuccessBold.Sprint(msg)
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
func (ui *basicUI) NamedValues(rows []NamedValue, opts ...Option) {
	cfg := &config{Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	cfg.Writer.Write([]byte{'\n'})

	var buf bytes.Buffer
	tr := tabwriter.NewWriter(&buf, 1, 8, 0, ' ', tabwriter.AlignRight)
	for _, row := range rows {
		fmt.Fprintf(tr, "  %s: \t%s\n", row.Name, row.Value)
	}

	tr.Flush()
	colorInfo.Fprintln(cfg.Writer, buf.String())
}

// OutputWriters implements UI
func (ui *basicUI) OutputWriters() (io.Writer, io.Writer, error) {
	return os.Stdout, os.Stderr, nil
}

// Status implements UI
func (ui *basicUI) Status() Status {
	if ui.status == nil {
		ui.status = newSpinnerStatus(ui.ctx)
	}

	return ui.status
}

func (ui *basicUI) StepGroup() StepGroup {
	ctx, cancel := context.WithCancel(ui.ctx)
	display := NewDisplay(ctx, color.Output)

	return &fancyStepGroup{
		ctx:     ctx,
		cancel:  cancel,
		display: display,
		done:    make(chan struct{}),
	}
}
