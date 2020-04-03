package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

// UI is the primary interface for interacting with a user via the CLI.
//
// NOTE(mitchellh): This is an interface and not a struct directly so that
// we can support other user interaction patterns in the future more easily.
// Most importantly what I'm thinking of is when we support multiple "apps"
// in a single config file, we can build a UI that locks properly and so on
// without changing the API.
type UI interface {
	// Output outputs a message directly to the terminal. The remaining
	// arguments should be interpolations for the format string. After the
	// interpolations you may add Options.
	Output(string, ...interface{})

	// OutputWriters returns stdout and stderr writers. These are usually
	// but not always TTYs. This is useful for subprocesses, network requests,
	// etc. Note that writing to these is not thread-safe by default so
	// you must take care that there is only ever one writer.
	OutputWriters() (stdout, stderr io.Writer, err error)

	// Status returns a live-updating status that can be used for single-line
	// status updates that typically have a spinner or some similar style.
	// While a Status is live (Close isn't called), Output should NOT be called.
	Status() Status
}

// BasicUI
type BasicUI struct{}

// Output implements UI
func (ui *BasicUI) Output(msg string, raw ...interface{}) {
	// Build our args and options
	var args []interface{}
	var opts []Option
	for _, r := range raw {
		if opt, ok := r.(Option); ok {
			opts = append(opts, opt)
		} else {
			args = append(args, r)
		}
	}

	// Build our message
	msg = fmt.Sprintf(msg, args...)

	// Build our config and set our options
	cfg := &config{Original: msg, Message: msg, Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	// Write it
	fmt.Fprintln(cfg.Writer, cfg.Message)
}

// OutputWriters implements UI
func (ui *BasicUI) OutputWriters() (io.Writer, io.Writer, error) {
	return os.Stdout, os.Stderr, nil
}

// Status implements UI
func (ui *BasicUI) Status() Status {
	return newSpinnerStatus()
}

type config struct {
	// Original is the original message, this should NOT be modified.
	Original string

	// Message is the message to write.
	Message string

	// Writer is where the message will be written to.
	Writer io.Writer
}

// Option controls output styling.
type Option func(*config)

// WithHeaderStyle styles the output like a header denoting a new section
// of execution. This should only be used with single-line output. Multi-line
// output will not look correct.
func WithHeaderStyle() Option {
	return func(c *config) {
		c.Message = colorHeader.Sprintf("==> %s", c.Message)
	}
}

// WithStatusStyle styles the output like a status update.
func WithStatusStyle() Option {
	return func(c *config) {
		lines := strings.Split(c.Message, "\n")
		for i, line := range lines {
			lines[i] = colorStatus.Sprintf("    %s", line)
		}

		c.Message = strings.Join(lines, "\n")
	}
}

// WithErrorStyle styles the output as an error message.
func WithErrorStyle() Option {
	return func(c *config) {
		c.Message = colorError.Sprint(c.Original)
	}
}

// WithWarningStyle styles the output as an error message.
func WithWarningStyle() Option {
	return func(c *config) {
		c.Message = colorWarning.Sprint(c.Original)
	}
}

// WithSuccessStyle styles the output as a success message.
func WithSuccessStyle() Option {
	return func(c *config) {
		c.Message = colorSuccess.Sprint(c.Original)
	}
}

// WithKeyValueStyle styles the output with aligned key/values with
// the given separator. This expects the the message is multiple lines
// which will be aligned. If a line doesn't contain a separator, it is
// ignored.
func WithKeyValueStyle(sep string) Option {
	return func(c *config) {
		// Trim whitespace first
		msg := strings.TrimSpace(c.Message)
		if len(msg) == 0 {
			return
		}

		// Go through each line, find the separator and record the whitespace.
		lines := strings.Split(msg, "\n")
		lineIdx := make([]int, len(lines))
		maxIdx := 0
		for i, line := range lines {
			lineIdx[i] = strings.Index(line, sep)
			if lineIdx[i] > maxIdx {
				maxIdx = lineIdx[i]
			}
		}

		// Output
		var buf bytes.Buffer
		for i, line := range lines {
			sepIdx := lineIdx[i]

			// Ignore lines with no sep
			if sepIdx < 0 {
				buf.WriteString(line)
				buf.WriteRune('\n')
				continue
			}

			// Pad
			buf.WriteString(strings.Repeat(" ", maxIdx-sepIdx))
			buf.WriteString(line)
			buf.WriteRune('\n')
		}

		bs := buf.Bytes()
		c.Message = string(bs[:len(bs)-1])
	}
}

// WithWriter specifies the writer for the output.
func WithWriter(w io.Writer) Option {
	return func(c *config) { c.Writer = w }
}

var (
	colorHeader  = color.New(color.Bold)
	colorStatus  = color.New()
	colorError   = color.New(color.FgRed)
	colorSuccess = color.New(color.FgGreen)
	colorWarning = color.New(color.FgYellow)
)
