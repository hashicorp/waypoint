package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/require"
)

func init() {
	goldie.FixtureDir = "testdata"
	spew.Config.DisablePointerAddresses = true
}

func TestParseFile(t *testing.T) {
	f, err := os.Open("testdata")
	require.NoError(t, err)
	defer f.Close()

	fis, err := f.Readdir(-1)
	require.NoError(t, err)
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if filepath.Ext(fi.Name()) == ".golden" {
			continue
		}

		t.Run(fi.Name(), func(t *testing.T) {
			var cfg Config
			err := hclsimple.DecodeFile(filepath.Join("testdata", fi.Name()), nil, &cfg)
			if strings.HasSuffix(fi.Name(), "_error.hcl") {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			goldie.Assert(t, fi.Name(), []byte(spew.Sdump(cfg)))
		})
	}
}
