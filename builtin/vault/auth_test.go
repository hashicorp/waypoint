// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthMethodConfig(t *testing.T) {
	cases := []struct {
		Name   string
		Config sourceConfig
		Result map[string]interface{}
	}{
		{
			"no method",
			sourceConfig{},
			map[string]interface{}{},
		},

		{
			"set values",
			sourceConfig{
				AuthMethod: "aws",
				AWSType:    "foo",
				AWSRole:    "bar",
				K8SRole:    "hello",
			},
			map[string]interface{}{
				"type": "foo",
				"role": "bar",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			result, _ := authMethodConfig(&tt.Config)
			require.Equal(result, tt.Result)
		})
	}
}
