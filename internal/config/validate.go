package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

// validateStruct is the validation structure for the configuration.
// This is used to validate the full structure of the configuration. This
// requires duplication between this struct and the other config structs
// since we don't do any lazy loading here.
type validateStruct struct {
	Project   string              `hcl:"project,optional"`
	Runner    *Runner             `hcl:"runner,block" default:"{}"`
	Labels    map[string]string   `hcl:"labels,optional"`
	Variables []*validateVariable `hcl:"variable,block"`
	Plugin    []*Plugin           `hcl:"plugin,block"`
	Apps      []*validateApp      `hcl:"app,block"`
	Pipelines []*validatePipeline `hcl:"pipeline,block"`
	Config    *genericConfig      `hcl:"config,block"`
}

type validateApp struct {
	Name    string            `hcl:",label"`
	Path    string            `hcl:"path,optional"`
	Labels  map[string]string `hcl:"labels,optional"`
	URL     *AppURL           `hcl:"url,block" default:"{}"`
	Runner  *Runner           `hcl:"runner,block"`
	Build   *Build            `hcl:"build,block"`
	Deploy  *Deploy           `hcl:"deploy,block"`
	Release *Release          `hcl:"release,block"`
	Config  *genericConfig    `hcl:"config,block"`
}

// validateVariable is separate from HclVariable because of the limitations
// of hclsimple, which we use in config.Load. hclsimple needs Type to be an
// hcl.Expression, but we want it to be a cty.Type for everything else.
type validateVariable struct {
	Name        string    `hcl:",label"`
	Default     cty.Value `hcl:"default,optional"`
	Type        cty.Type  `hcl:"type,optional"`
	Description string    `hcl:"description,optional"`
}

type validatePipeline struct {
	Name   string            `hcl:",label"`
	Labels map[string]string `hcl:"labels,optional"`
	Step   []*Step           `hcl:"step,block"`
}

type ValidationResult struct {
	Error   error
	Warning string
}

func (v ValidationResult) String() string {
	if v.Error != nil {
		return v.Error.Error()
	}

	return "warning: " + v.Warning
}

type ValidationResults []ValidationResult

func (v ValidationResults) Error() string {
	var values []string

	for _, res := range v {
		values = append(values, res.String())
	}

	return fmt.Sprintf("%d validation errors: %s", len(v), strings.Join(values, ", "))
}

func (v ValidationResults) HasErrors() bool {
	for _, vr := range v {
		if vr.Error != nil {
			return true
		}
	}

	return false
}

const AppeningHappening = "[NOTICE] More than one app stanza within a waypoint.hcl file is under consideration for change or removal in a future version.\nTo give feedback, visit https://discuss.hashicorp.com/t/deprecating-projects-or-how-i-learned-to-love-apps/40888"

// Validate the structure of the configuration.
//
// This will validate required fields are specified and the types of some fields.
// Plugin-specific fields won't be validated until later. Fields that use functions
// and variables will not be validated until those values can be realized.
//
// Users of this package should call Validate on each subsequent configuration
// that is loaded (Apps, Builds, Deploys, etc.) for further rich validation.
func (c *Config) Validate() (ValidationResults, error) {
	var results ValidationResults

	// Validate root
	schema, _ := gohcl.ImpliedBodySchema(&validateStruct{})
	content, diag := c.hclConfig.Body.Content(schema)
	if diag.HasErrors() {
		results = append(results, ValidationResult{Error: diag})

		if content == nil {
			return results, results
		}
	}

	// Require the project. We don't use an "attr" above (which would require it)
	// because the project can be populated later such as in a runner which
	// sets it to the project in the job ref.
	if c.Project == "" {
		results = append(results, ValidationResult{Error: fmt.Errorf("'project' attribute is required")})
	}

	apps := content.Blocks.OfType("app")

	// Validate apps
	for _, block := range apps {
		appRes := c.validateApp(block)
		results = append(results, appRes...)
	}

	if len(apps) > 1 {
		results = append(results, ValidationResult{Warning: AppeningHappening})
	}

	// Validate pipelines
	for i, block := range content.Blocks.OfType("pipeline") {
		pipelineRes := c.validatePipeline(block)
		if pipelineRes != nil {
			results = append(results, pipelineRes...)
		}

		// Validate there's no duplicate names
		for j, bl := range content.Blocks.OfType("pipeline") {
			if i != j && bl.Labels[0] == block.Labels[0] {
				results = append(results, ValidationResult{Error: &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "'pipeline' stanza names must be unique per project",
					Subject:  &block.DefRange,
					Context:  &block.TypeRange,
				}})
			}
		}
	}

	// Validate labels
	labelResults := ValidateLabels(c.Labels)
	results = append(results, labelResults...)

	// So that callers that test the result for nil can still do so
	// (they test for nil because this return type used to be error)
	if len(results) == 0 {
		return results, nil
	}

	if results.HasErrors() {
		return results, results
	}

	return results, nil
}

func (c *Config) validateApp(b *hcl.Block) []ValidationResult {
	var results []ValidationResult

	// Validate root
	schema, _ := gohcl.ImpliedBodySchema(&validateApp{})
	content, diag := b.Body.Content(schema)
	if diag.HasErrors() {
		results = append(results, ValidationResult{Error: diag})

		if content == nil {
			return results
		}
	}

	// Build required
	if len(content.Blocks.OfType("build")) != 1 {
		results = append(results, ValidationResult{Error: &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "'build' stanza required",
			Subject:  &b.DefRange,
			Context:  &b.TypeRange,
		}})
	}

	// Deploy required
	if len(content.Blocks.OfType("deploy")) != 1 {
		results = append(results, ValidationResult{Error: &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "'deploy' stanza required",
			Subject:  &b.DefRange,
			Context:  &b.TypeRange,
		}})
	}

	return results
}

// Validate validates the application.
//
// Similar to Config.App, this doesn't validate configuration that is
// further deferred such as build, deploy, etc. stanzas so call Validate
// on those as they're loaded.
func (c *App) Validate() (ValidationResults, error) {
	var results ValidationResults

	// Validate labels
	labelResults := ValidateLabels(c.Labels)
	results = append(results, labelResults...)

	// If a path is specified, it must not be a child of the root.
	if c.Path != "" {
		if !filepath.IsAbs(c.Path) {
			// This should never happen because during App load time
			// we ensure that the path is absolute relative to the project
			// path.
			results = append(results, ValidationResult{Error: fmt.Errorf("path is not absolute")})
		}

		rel, err := filepath.Rel(c.config.path, c.Path)
		if err != nil {
			results = append(results, ValidationResult{Error: fmt.Errorf(
				"path: must be a child of the project directory")})
		}
		if strings.HasPrefix(rel, "../") || strings.HasPrefix(rel, "..\\") {
			results = append(results, ValidationResult{Error: fmt.Errorf(
				"path: must be a child of the project directory")})
		}
	}

	if c.BuildRaw == nil || c.BuildRaw.Use == nil || c.BuildRaw.Use.Type == "" {
		results = append(results, ValidationResult{Error: fmt.Errorf(
			"build stage with a default non-workspace scoped 'use' stanza is required")})
	}

	if c.DeployRaw == nil || c.DeployRaw.Use == nil || c.DeployRaw.Use.Type == "" {
		results = append(results, ValidationResult{Error: fmt.Errorf(
			"deploy stage with a default non-workspace scoped 'use' stanza is required")})
	}

	for _, scope := range c.BuildRaw.WorkspaceScoped {
		if scope.Use == nil || scope.Use.Type == "" {
			results = append(results, ValidationResult{Error: fmt.Errorf(
				"build: workspace scope %q: 'use' stanza is required",
				scope.Scope,
			)})
		}
	}
	for _, scope := range c.BuildRaw.LabelScoped {
		if scope.Use == nil || scope.Use.Type == "" {
			results = append(results, ValidationResult{Error: fmt.Errorf(
				"build: label scope %q: 'use' stanza is required",
				scope.Scope,
			)})
		}
	}

	for _, scope := range c.DeployRaw.WorkspaceScoped {
		if scope.Use == nil || scope.Use.Type == "" {
			results = append(results, ValidationResult{Error: fmt.Errorf(
				"deploy: workspace scope %q: 'use' stanza is required",
				scope.Scope,
			)})
		}
	}
	for _, scope := range c.DeployRaw.LabelScoped {
		if scope.Use == nil || scope.Use.Type == "" {
			results = append(results, ValidationResult{Error: fmt.Errorf(
				"deploy: label scope %q: 'use' stanza is required",
				scope.Scope,
			)})
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	if results.HasErrors() {
		return results, results
	}

	return results, nil
}

// validatePipeline validates that a given pipeline block has at least
// one step stanza
func (c *Config) validatePipeline(b *hcl.Block) []ValidationResult {
	var results []ValidationResult

	// Validate root
	schema, _ := gohcl.ImpliedBodySchema(&validatePipeline{})
	content, diag := b.Body.Content(schema)
	if diag.HasErrors() {
		results = append(results, ValidationResult{Error: diag})

		if content == nil {
			return results
		}
	}

	// At least one Step required
	if len(content.Blocks.OfType("step")) < 1 {
		results = append(results, ValidationResult{Error: &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "'step' stanza required",
			Subject:  &b.DefRange,
			Context:  &b.TypeRange,
		}})
	}

	return results
}

func (c *Pipeline) Validate() error {
	var result error

	for _, stepRaw := range c.StepRaw {
		if stepRaw == nil && stepRaw.PipelineRaw == nil {
			result = multierror.Append(result, fmt.Errorf(
				"step stage with a default 'use' stanza or a 'pipeline' stanza is required"))
		} else if stepRaw.Use != nil && stepRaw.PipelineRaw != nil {
			result = multierror.Append(result, fmt.Errorf(
				"step stage with a 'use' stanza and pipeline stanza is not valid"))
		} else if stepRaw.PipelineRaw == nil && (stepRaw.Use == nil || stepRaw.Use.Type == "") {
			result = multierror.Append(result, fmt.Errorf(
				"step stage %q is required to define a 'use' stanza and label", stepRaw.Name))
		}

		// else, other step validations?
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
func ValidateLabels(labels map[string]string) ValidationResults {
	var results ValidationResults

	for k, v := range labels {
		name := fmt.Sprintf("label[%s]", k)

		if strings.HasPrefix(k, "waypoint/") {
			results = append(results, ValidationResult{Error: fmt.Errorf("%s: prefix 'waypoint/' is reserved for system use", name)})
		}

		if len(k) > 255 {
			results = append(results, ValidationResult{Error: fmt.Errorf("%s: key must be less than or equal to 255 characters", name)})
		}

		if !hostnameRegexRFC952.MatchString(strings.SplitN(k, "/", 2)[0]) {
			results = append(results, ValidationResult{Error: fmt.Errorf("%s: key before '/' must be a valid hostname (RFC 952)", name)})
		}

		if len(v) > 255 {
			results = append(results, ValidationResult{Error: fmt.Errorf("%s: value must be less than or equal to 255 characters", name)})
		}
	}

	return results
}

var hostnameRegexRFC952 = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`)
