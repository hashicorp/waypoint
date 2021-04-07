package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// validateStruct is the validation structure for the configuration.
// This is used to validate the full structure of the configuration. This
// requires duplication between this struct and the other config structs
// since we don't do any lazy loading here.
type validateStruct struct {
	Project string            `hcl:"project,optional"`
	Runner  *Runner           `hcl:"runner,block" default:"{}"`
	Labels  map[string]string `hcl:"labels,optional"`
	Plugin  []*Plugin         `hcl:"plugin,block"`
	Apps    []*validateApp    `hcl:"app,block"`
	Config  *genericConfig    `hcl:"config,block"`
}

type validateApp struct {
	Name    string            `hcl:",label"`
	Path    string            `hcl:"path,optional"`
	Labels  map[string]string `hcl:"labels,optional"`
	URL     *AppURL           `hcl:"url,block" default:"{}"`
	Build   *Build            `hcl:"build,block"`
	Deploy  *Deploy           `hcl:"deploy,block"`
	Release *Release          `hcl:"release,block"`
	Config  *genericConfig    `hcl:"config,block"`
}

// Validate the structure of the configuration.
//
// This will validate required fields are specified and the types of some fields.
// Plugin-specific fields won't be validated until later. Fields that use functions
// and variables will not be validated until those values can be realized.
//
// Users of this package should call Validate on each subsequent configuration
// that is loaded (Apps, Builds, Deploys, etc.) for further rich validation.
func (c *Config) Validate() error {
	// Validate root
	schema, _ := gohcl.ImpliedBodySchema(&validateStruct{})
	content, diag := c.hclConfig.Body.Content(schema)
	if diag.HasErrors() {
		return diag
	}

	var result error

	// Require the project. We don't use an "attr" above (which would require it)
	// because the project can be populated later such as in a runner which
	// sets it to the project in the job ref.
	if c.Project == "" {
		result = multierror.Append(result, fmt.Errorf("'project' attribute is required"))
	}

	// Validate apps
	for _, block := range content.Blocks.OfType("app") {
		err := c.validateApp(block)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Validate labels
	if errs := ValidateLabels(c.Labels); len(errs) > 0 {
		result = multierror.Append(result, errs...)
	}

	return result
}

func (c *Config) validateApp(b *hcl.Block) error {
	// Validate root
	schema, _ := gohcl.ImpliedBodySchema(&validateApp{})
	content, diag := b.Body.Content(schema)
	if diag.HasErrors() {
		return diag
	}

	// Build required
	if len(content.Blocks.OfType("build")) != 1 {
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "'build' stanza required",
			Subject:  &b.DefRange,
			Context:  &b.TypeRange,
		}
	}

	// Deploy required
	if len(content.Blocks.OfType("deploy")) != 1 {
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "'deploy' stanza required",
			Subject:  &b.DefRange,
			Context:  &b.TypeRange,
		}
	}

	return nil
}

// Validate validates the application.
//
// Similar to Config.App, this doesn't validate configuration that is
// further deferred such as build, deploy, etc. stanzas so call Validate
// on those as they're loaded.
func (c *App) Validate() error {
	var result error

	// Validate labels
	if errs := ValidateLabels(c.Labels); len(errs) > 0 {
		result = multierror.Append(result, errs...)
	}

	// If a path is specified, it must not be a child of the root.
	if c.Path != "" {
		if !filepath.IsAbs(c.Path) {
			// This should never happen because during App load time
			// we ensure that the path is absolute relative to the project
			// path.
			panic("path is not absolute")
		}

		rel, err := filepath.Rel(c.config.path, c.Path)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf(
				"path: must be a child of the project directory"))
		}
		if strings.HasPrefix(rel, "../") || strings.HasPrefix(rel, "..\\") {
			result = multierror.Append(result, fmt.Errorf(
				"path: must be a child of the project directory"))
		}
	}

	if c.BuildRaw == nil || c.BuildRaw.Use == nil || c.BuildRaw.Use.Type == "" {
		result = multierror.Append(result, fmt.Errorf(
			"build stage with a 'use' stanza is required"))
	}

	if c.DeployRaw == nil || c.DeployRaw.Use == nil || c.DeployRaw.Use.Type == "" {
		result = multierror.Append(result, fmt.Errorf(
			"deploy stage with a 'use' stanza is required"))
	}

	return result
}

// ValidateLabels validates a set of labels. This ensures that labels are
// set according to our requirements:
//
//   * key and value length can't be greater than 255 characters each
//   * keys must be in hostname format (RFC 952)
//   * keys can't be prefixed with "waypoint/" which is reserved for system use
//
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
