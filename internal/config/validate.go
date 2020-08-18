package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// internalValidator is an interface implemented internally for validation.
// This is unexported since it takes a name parameter to build better error
// messages.
type internalValidator interface {
	validate(name string) error
}

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

	// If a path is specified, it must be relative and it must not
	// contain any "..".
	if app.Path != "" {
		if filepath.IsAbs(app.Path) {
			result = multierror.Append(result, fmt.Errorf(
				"path: must be a relative path"))
		}

		for _, part := range filepath.SplitList(app.Path) {
			if part == ".." {
				result = multierror.Append(result, fmt.Errorf(
					"path: must not contain .. entries"))
				break
			}
		}
	}

	for k, v := range app.validatorChildren() {
		if v != nil {
			if err := v.validate(k); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	return multierror.Prefix(result, fmt.Sprintf("app[%s]:", app.Name))
}

func (app *App) validatorChildren() map[string]internalValidator {
	result := map[string]internalValidator{
		"build":   app.Build,
		"deploy":  app.Deploy,
		"release": app.Release,
	}

	if app.Build != nil && app.Build.Registry != nil {
		result["build.registry"] = app.Build.Registry
	}

	return result
}

func (c *Build) validate(key string) error {
	return c.Operation().validate(key)
}

func (c *Deploy) validate(key string) error {
	return c.Operation().validate(key)
}

func (c *Registry) validate(key string) error {
	return c.Operation().validate(key)
}

func (c *Release) validate(key string) error {
	return c.Operation().validate(key)
}

func (c *Operation) validate(key string) error {
	if c == nil {
		return nil
	}

	var result error
	if errs := ValidateLabels(c.Labels); len(errs) > 0 {
		result = multierror.Append(result, errs...)
	}

	for i, h := range c.Hooks {
		if err := h.validate(fmt.Sprintf("hook[%d]", i)); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return multierror.Prefix(result, fmt.Sprintf("%s:", key))
}

func (h *Hook) validate(key string) error {
	var result error

	switch h.When {
	case "before", "after":
	default:
		result = multierror.Append(result, fmt.Errorf("label must be 'before' or 'after'"))
	}

	if len(h.Command) == 0 {
		result = multierror.Append(result, fmt.Errorf("command must be non-empty"))
	}

	switch h.OnFailure {
	case "", "continue", "fail":
	default:
		result = multierror.Append(result, fmt.Errorf("on_failure must be 'continue' or 'fail'"))
	}

	return multierror.Prefix(result, fmt.Sprintf("%s:", key))
}

// ValidateLabels validates a set of labels.
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
