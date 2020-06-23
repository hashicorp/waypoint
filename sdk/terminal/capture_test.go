package terminal

import (
	"io"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func TestCaptureUI_output(t *testing.T) {
	require := require.New(t)

	// Create our UI
	total := ""
	ui := NewCaptureUI(hclog.L(), func(lines []*CaptureLine) error {
		for _, line := range lines {
			total += line.Line + "\n"
		}

		return nil
	})
	defer ui.Close()

	// Simple output
	ui.Output("hello")
	ui.Output("there")
	ui.Output("world")
	ui.Close()

	// Verify
	require.Equal(`hello
there
world
`, total)
}

func TestCaptureUI_outputWriter(t *testing.T) {
	require := require.New(t)

	// Create our UI
	total := ""
	ui := NewCaptureUI(hclog.L(), func(lines []*CaptureLine) error {
		for _, line := range lines {
			total += line.Line + "\n"
		}

		return nil
	})
	defer ui.Close()

	// Get our writer
	out, _, err := ui.OutputWriters()
	require.NoError(err)

	// Write to it
	_, err = io.Copy(out, strings.NewReader(`hello
another line`))
	require.NoError(err)

	// Close it up
	ui.Close()

	// Verify
	require.Equal(`hello
another line
`, total)
}
