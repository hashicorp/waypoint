// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import "os"

func ConsulAuth() (string, error) {
	return os.Getenv("CONSUL_HTTP_TOKEN"), nil
}

func VaultAuth() (string, error) {
	return os.Getenv("VAULT_TOKEN"), nil
}
