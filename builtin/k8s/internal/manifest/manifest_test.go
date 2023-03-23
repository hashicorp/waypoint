// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	f, err := os.Open(filepath.Join("testdata", "deployment.yaml"))
	require.NoError(err)
	defer f.Close()

	m, err := Parse(f)
	require.NoError(err)
	require.Len(m.Resources, 2)

	{
		r := m.Resources[0]
		require.Equal("Deployment", r.Kind)
		require.Equal("php-apache", r.Metadata.Name)
	}
}
