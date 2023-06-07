// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helm

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-hclog"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	memcached "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashicorp/waypoint/builtin/k8s"
)

// restClientGetter is a RESTClientGetter interface implementation for the
// Helm Go packages.
type restClientGetter struct {
	RestConfig  *rest.Config
	Kubeconfig  string
	Kubecontext string
}

// ToRESTConfig implemented interface method
func (k *restClientGetter) ToRESTConfig() (*rest.Config, error) {
	return k.RestConfig, nil
}

// ToDiscoveryClient implemented interface method
func (k *restClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := k.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom resources) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	config.Burst = 100

	return memcached.NewMemCacheClient(discovery.NewDiscoveryClientForConfigOrDie(config)), nil
}

// ToRESTMapper implemented interface method
func (k *restClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := k.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

// ToRawKubeConfigLoader implemented interface method
func (k *restClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	loader := clientcmd.NewDefaultClientConfigLoadingRules()

	// Path to the kube config file
	if k.Kubeconfig != "" {
		loader.ExplicitPath = k.Kubeconfig
	}

	// Build our config and client
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loader,
		&clientcmd.ConfigOverrides{
			CurrentContext: k.Kubecontext,
		},
	)

	return config
}

func SettingsInit(ns string) (*cli.EnvSettings, error) {
	cli := cli.New()
	if ns != "" {
		cli.SetNamespace(ns)
	}
	return cli, nil
}

func ActionInit(log hclog.Logger, kubeConfigPath string, context string, namespace string) (*action.Configuration, error) {
	// Get our K8S API
	_, ns, rc, err := k8s.Clientset(kubeConfigPath, context)
	if err != nil {
		return nil, err
	}
	if namespace != "" {
		ns = namespace
	}
	driver := "secret"

	// For logging, we'll debug log to a custom named logger.
	actionlog := log.Named("helm_action")
	debug := func(format string, v ...interface{}) {
		actionlog.Debug(fmt.Sprintf(format, v...))
	}

	// Initialize our action
	var ac action.Configuration
	err = ac.Init(&restClientGetter{
		RestConfig:  rc,
		Kubeconfig:  kubeConfigPath,
		Kubecontext: context,
	}, ns, driver, debug)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}

func ChartPathOptions(repository string, chart string, version string) (*action.ChartPathOptions, string, error) {
	repositoryURL, chartName, err := resolveChartName(
		repository, strings.TrimSpace(chart))
	if err != nil {
		return nil, "", err
	}

	// Determine our version string
	version = strings.TrimSpace(version)

	// Initialize our chart options
	return &action.ChartPathOptions{
		RepoURL: repositoryURL,
		Version: version,
	}, chartName, nil
}

func GetChart(name string, cpo *action.ChartPathOptions, settings *cli.EnvSettings) (*chart.Chart, string, error) {
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

	if !strings.Contains(name, "/") && repository != "" {
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

func GetLatestHelmChartVersion(ctx context.Context) ([]*github.RepositoryTag, error) {
	githubClient := github.NewClient(nil)
	tags, _, err := githubClient.Repositories.ListTags(ctx, "hashicorp", "waypoint-helm", nil)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
