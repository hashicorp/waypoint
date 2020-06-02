package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func (c *Config) Validate() error {
	var result error

	if errs := ValidateLabels(c.Labels); len(errs) > 0 {
		result = multierror.Append(result, errs...)
	}

	for _, app := range c.Apps {
		if err := app.Validate(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func (app *App) Validate() error {
	var result error
	if errs := ValidateLabels(app.Labels); len(errs) > 0 {
		result = multierror.Append(result, errs...)
	}

	for k, v := range app.components() {
		if v != nil {
			if err := v.validate(k); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	return multierror.Prefix(result, fmt.Sprintf("app[%s]:", app.Name))
}

func (c *Component) validate(key string) error {
	var result error
	if errs := ValidateLabels(c.Labels); len(errs) > 0 {
		result = multierror.Append(result, errs...)
	}

	return multierror.Prefix(result, fmt.Sprintf("%s:", key))
}

func ValidateLabels(labels map[string]string) []error {
	var errs []error
	for k, v := range labels {
		name := fmt.Sprintf("label[%s]", k)

		if strings.HasPrefix(k, "waypoint/") {
			errs = append(errs, fmt.Errorf("%s: prefix 'waypoint/' is reserved for system use", name))
		}

		if len(k) > 255 {
			errs = append(errs, fmt.Errorf("%s: key must be less than or equal to 255 characters", name))
		}

		if !hostnameRegexRFC952.MatchString(strings.SplitN(k, "/", 2)[0]) {
			errs = append(errs, fmt.Errorf("%s: key before '/' must be a valid hostname (RFC 952)", name))
		}

		if len(v) > 255 {
			errs = append(errs, fmt.Errorf("%s: value must be less than or equal to 255 characters", name))
		}
	}

	return errs
}

var hostnameRegexRFC952 = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`)
