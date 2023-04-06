// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validationext

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	// IsDuration implements validation.Rule to check if a value is a valid
	// duration from a proto duration.
	IsDuration validation.Rule = &isDurationRule{}
)

// IsDurationRange implements validation.Rule to check if a proto duration
// is in the given range. Both ends of the duration are inclusive.
func IsDurationRange(min, max time.Duration) validation.Rule {
	return &isDurationRange{min: min, max: max}
}

// isDurationRule implements validation.Rule for IsDuration
type isDurationRule struct{}

func (r *isDurationRule) Validate(value interface{}) error {
	_, err := r.duration(value)
	return err
}

func (r *isDurationRule) duration(value interface{}) (time.Duration, error) {
	switch v := value.(type) {
	case *durationpb.Duration:
		// Support non-required duration.
		if v == nil {
			return 0, nil
		}
		return v.AsDuration(), nil
	case string:
		// Support non-required duration.
		if v == "" {
			return 0, nil
		}

		return time.ParseDuration(v)
	case *time.Duration:
		// Support non-required duration.
		if v == nil {
			return 0, nil
		}

		return *v, nil
	case time.Duration:
		return v, nil
	}

	return 0, fmt.Errorf("must be a valid duration")
}

// isDurationRange implements validation.Rule for IsDurationRange.
type isDurationRange struct {
	min, max time.Duration
}

func (r *isDurationRange) Validate(value interface{}) error {
	var dr isDurationRule
	d, err := dr.duration(value)
	if err != nil {
		return err
	}
	if d < r.min || d > r.max {
		return fmt.Errorf(
			"must be greater than %s and less than %s",
			r.min.String(), r.max.String())
	}

	return nil
}
