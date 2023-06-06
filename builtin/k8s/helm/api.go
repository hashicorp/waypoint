package helm

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"

	"github.com/hashicorp/waypoint/builtin/k8s"
)

func (p *Platform) settingsInit() (*cli.EnvSettings, error) {
	cli := cli.New()
	if p.config.Namespace != "" {
		cli.SetNamespace(p.config.Namespace)
	}
	return cli, nil
}

func (p *Platform) actionInit(log hclog.Logger) (*action.Configuration, error) {
	// Get our K8S API
	_, ns, rc, err := k8s.Clientset(p.config.KubeconfigPath, p.config.Context)
	if err != nil {
		return nil, err
	}

	if p.config.Namespace != "" {
		ns = p.config.Namespace
	}

	driver := "secret"
	if v := p.config.Driver; v != "" {
		driver = v
	}

	// For logging, we'll debug log to a custom named logger.
	actionlog := log.Named("helm_action")
	debug := func(format string, v ...interface{}) {
		actionlog.Debug(fmt.Sprintf(format, v...))
	}

	// Initialize our action
	var ac action.Configuration
	err = ac.Init(&restClientGetter{
		RestConfig:  rc,
		Kubeconfig:  p.config.KubeconfigPath,
		Kubecontext: p.config.Context,
	}, ns, driver, debug)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}

func (p *Platform) chartPathOptions() (*action.ChartPathOptions, string, error) {
	repositoryURL, chartName, err := resolveChartName(
		p.config.Repository, strings.TrimSpace(p.config.Chart))
	if err != nil {
		return nil, "", err
	}

	// Determine our version string
	version := p.config.Version
	if version == "" && p.config.Devel {
		version = ">0.0.0-0"
	}
	version = strings.TrimSpace(version)

	// Initialize our chart options
	return &action.ChartPathOptions{
		RepoURL: repositoryURL,
		Version: version,
	}, chartName, nil
}

func (p *Platform) chartValues() (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// First merge all our values from YAML documents.
	for _, values := range p.config.Values {
		if values == "" {
			continue
		}

		currentMap := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(values), &currentMap); err != nil {
			return nil, fmt.Errorf("---> %v %s", err, values)
		}

		base = mergeMaps(base, currentMap)
	}

	// Next get all our set configs
	for _, set := range p.config.Set {
		name := set.Name
		value := set.Value
		valueType := set.Type

		switch valueType {
		case "auto", "":
			if err := strvals.ParseInto(fmt.Sprintf("%s=%s", name, value), base); err != nil {
				return nil, fmt.Errorf("failed parsing key %q with value %s, %s", name, value, err)
			}
		case "string":
			if err := strvals.ParseIntoString(fmt.Sprintf("%s=%s", name, value), base); err != nil {
				return nil, fmt.Errorf("failed parsing key %q with value %s, %s", name, value, err)
			}
		default:
			return nil, fmt.Errorf("unexpected type: %s", valueType)
		}
	}

	return base, nil
}

func getChart(name string, cpo *action.ChartPathOptions, settings *cli.EnvSettings) (*chart.Chart, string, error) {
	path, err := cpo.LocateChart(name, settings)
	if err != nil {
		return nil, "", err
	}

	c, err := loader.Load(path)
	if err != nil {
		return nil, "", err
	}

	return c, path, nil
}

func getRelease(cfg *action.Configuration, name string) (*release.Release, error) {
	res, err := action.NewGet(cfg).Run(name)
	if err != nil {
		if strings.Contains(err.Error(), "release: not found") {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

// resolveChartName returns the proper repository and name values that
// the ChartPathOptions need. This is copied from Terraform.
func resolveChartName(repository, name string) (string, string, error) {
	_, err := url.ParseRequestURI(repository)
	if err == nil {
		return repository, name, nil
	}

	if strings.Index(name, "/") == -1 && repository != "" {
		name = fmt.Sprintf("%s/%s", repository, name)
	}

	return "", name, nil
}

// Merges source and destination map, preferring values from the source map
// Taken from github.com/helm/pkg/cli/values/options.go
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
