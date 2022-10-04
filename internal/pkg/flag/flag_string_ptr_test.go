package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPtrVar(t *testing.T) {
	cases := map[string]struct {
		input     *string
		expected  *string
		omitValue bool // e.g. -poll instead of -poll=true
		shouldErr bool
	}{
		"omitted": {
			expected: nil,
		},
		// specifying -flag without a value should fail
		"no value set": {
			omitValue: true,
			shouldErr: true,
		},
		"value passed in": {
			input:    strPtr("blammo"),
			expected: strPtr("blammo"),
		},
		// empty is equivilant to a user supplying -flag="", and not a scenario
		// where the flag is simply omittied.
		"empty string as input": {
			input:    strPtr(""),
			expected: strPtr(""),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var valA *string
			var valB string
			sets := NewSets()
			{
				set := sets.NewSet("A")
				set.StringPtrVar(&StringPtrVar{
					Name:   "a",
					Target: &valA,
				})
			}
			{
				// borrowed from string_slice_test, just to have another input
				set := sets.NewSet("B")
				set.StringVar(&StringVar{
					Name:   "b",
					Target: &valB,
				})
			}

			var err error
			inputs := []string{"-b=somevalueB"}
			if c.input != nil {
				inputs = append(inputs, "-a="+*c.input)
			}
			if c.omitValue == true {
				inputs = append(inputs, "-a")
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
