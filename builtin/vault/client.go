// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"github.com/hashicorp/go-hclog"
	vaultapi "github.com/hashicorp/vault/api"
)

// initClient setes the cs.client value. If that value is not nil already, then
// this does nothing.
//
// This function expects cacheMu mutex is locked already.
func (cs *ConfigSourcer) initClient(log hclog.Logger) error {
	if cs.client != nil {
		return nil
	}

	// cs.Client is used for testing, use that if set
	if cs.Client != nil {
		log.Debug("using provided Client")
		cs.client = cs.Client
		return nil
	}

	// Start with the default config, then layer on any source config.
	// The default config reads env vars.
	clientConfig := vaultapi.DefaultConfig()
	if v := cs.config.Address; v != "" {
		clientConfig.Address = v
	}
	if v := cs.config.AgentAddress; v != "" {
		clientConfig.AgentAddress = v
	}
	var tlsConfig vaultapi.TLSConfig
	if v := cs.config.CACert; v != "" {
		tlsConfig.CACert = v
	}
	if v := cs.config.CAPath; v != "" {
		tlsConfig.CAPath = v
	}
	if v := cs.config.ClientCert; v != "" {
		tlsConfig.ClientCert = v
	}
	if v := cs.config.ClientKey; v != "" {
		tlsConfig.ClientKey = v
	}
	if v := cs.config.TLSServerName; v != "" {
		tlsConfig.TLSServerName = v
	}
	if v := cs.config.SkipVerify; v {
		tlsConfig.Insecure = v
	}
	if err := clientConfig.ConfigureTLS(&tlsConfig); err != nil {
		return err
	}

	log.Debug("initializing the Vault client")
	client, err := vaultapi.NewClient(clientConfig)
	if err != nil {
		return err
	}

	// Additional config vars that are set on the client
	if v := cs.config.Namespace; v != "" {
		client.SetNamespace(v)
	}
	if v := cs.config.Token; v != "" {
		client.SetToken(v)
	}

	cs.client = client
	return nil
}
