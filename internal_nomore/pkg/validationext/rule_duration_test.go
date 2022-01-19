package validationext

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
)

func TestIsDuration(t *testing.T) {
	cases := []struct {
		Input interface{}
		Valid bool
	}{
		{
			nil,
			false,
		},

		{
			"foo",
			false,
		},

		{
			ptypes.DurationProto(time.Second),
			true,
		},

		{
			"2s",
			true,
		},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%#v", tt.Input), func(t *testing.T) {
			err := IsDuration.Validate(tt.Input)
			require.Equal(t, tt.Valid, err == nil)
		})
	}
}

func TestIsDurationRange(t *testing.T) {
	cases := []struct {
		Duration interface{}
		Min, Max time.Duration
		Valid    bool
	}{
		{
			nil,
			0, 0,
			false,
		},

		{
			"foo",
			0, 0,
			false,
		},

		{
			ptypes.DurationProto(time.Second),
			0, time.Minute,
			true,
		},

		{
			ptypes.DurationProto(time.Minute),
			0, time.Second,
			false,
		},

		{
			ptypes.DurationProto(time.Second),
			time.Second, time.Minute,
			true,
		},

		{
			ptypes.DurationProto(time.Minute),
			time.Second, time.Minute,
			true,
		},

		{
			"1m",
			0, time.Second,
			false,
		},

		{
			"1s",
			time.Second, time.Minute,
			true,
		},

		{
			"1m",
			time.Second, time.Minute,
			true,
		},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%#v", tt.Duration), func(t *testing.T) {
			r := IsDurationRange(tt.Min, tt.Max)
			err := r.Validate(tt.Duration)
			require.Equal(t, tt.Valid, err == nil)
		})
	}
}
