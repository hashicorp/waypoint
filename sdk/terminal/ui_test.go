package terminal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	require := require.New(t)

	var buf bytes.Buffer
	var ui BasicUI
	ui.NamedValues([]NamedValue{
		{"hello", "a"},
		{"this", "is"},
		{"a", "test"},
		{"of", "foo"},
		{"the_key_value", "style"},
	},
		WithWriter(&buf),
	)

	expected := `
        hello: a
         this: is
            a: test
           of: foo
the_key_value: style

`

	require.Equal(expected, buf.String())
}

func TestStatusStyle(t *testing.T) {
	require := require.New(t)

	var buf bytes.Buffer
	var ui BasicUI
	ui.Output(strings.TrimSpace(`
one
two
  three`),
		WithWriter(&buf),
		WithInfoStyle(),
	)

	expected := `    one
    two
      three
`

	require.Equal(expected, buf.String())
}
