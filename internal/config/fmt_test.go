// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/config"
)

func TestFormat(t *testing.T) {
	const outSuffix = ".out"
	path := filepath.Join("testdata", "fmt")
	entries, err := ioutil.ReadDir(path)
	require.NoError(t, err)

	g := goldie.New(t,
		goldie.WithFixtureDir(filepath.Join("testdata", "fmt")),
		goldie.WithNameSuffix(outSuffix),
	)

	for _, entry := range entries {
		// Ignore golden files
		if strings.HasSuffix(entry.Name(), outSuffix) {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			require := require.New(t)

			// Read the input file
			src, err := ioutil.ReadFile(filepath.Join(path, entry.Name()))
			require.NoError(err)

			// Format it!
			out, err := config.Format(src, entry.Name())
			require.NoError(err)

			// Compare
			g.Assert(t, entry.Name(), out)
		})
	}
}
