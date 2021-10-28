package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringBool(t *testing.T) {
	cases := map[string]struct {
		input     *string
		expected  string
		shouldErr bool
	}{
		"omitted": {},
		"true": {
			input:    strPtr("TRUE"),
			expected: "true",
		},
		"false": {
			input:    strPtr("False"),
			expected: "false",
		},
		// empty is equivilant to a user supplying -flag="", and not a scenario
		// where the flag is simply omittied.
		"empty": {
			input:     strPtr(""),
			shouldErr: true,
		},
		"bad input": {
			input:     strPtr("a non-truthy value"),
			shouldErr: true,
		},
		// only accept values Go's ParseBool would accept
		"weird truthy input": {
			input:     strPtr("tRuE"),
			shouldErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var valA string
			var valB []string
			sets := NewSets()
			{
				set := sets.NewSet("A")
				set.StringBoolVar(&StringBoolVar{
					Name:   "a",
					Target: &valA,
				})
			}
			{
				// borrowed from string_slice_test, just to have another input
				set := sets.NewSet("B")
				set.StringSliceVar(&StringSliceVar{
					Name:   "b",
					Target: &valB,
				})
			}

			var err error
			inputs := []string{"-b", "somevalueB"}
			if c.input != nil {
				inputs = append(inputs, "-a", *c.input)
			}
			err = sets.Parse(inputs)
			if c.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, c.expected, valA)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
