// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helm

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	memcached "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
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
