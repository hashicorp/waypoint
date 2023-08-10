// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ec2

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func UserData(env map[string]string) (string, error) {
	envData, err := json.Marshal(env)
	if err != nil {
		return "", err
	}

	envStr := base64.StdEncoding.EncodeToString(envData)

	template := fmt.Sprintf(strings.TrimSpace(`
#cloud-config
write_files:
- encoding: b64
  content: %s
  owner: root:root
  path: /etc/waypoint/env
  permissions: '0644'
`), envStr)

	return base64.StdEncoding.EncodeToString([]byte(template)), nil
}
