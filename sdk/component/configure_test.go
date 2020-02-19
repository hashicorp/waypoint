package component

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/stretchr/testify/require"
)

func TestConfigure(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		require := require.New(t)

		var c impl
		src := `name = "foo"`
		f, diag := hclparse.NewParser().ParseHCL([]byte(src), "test.hcl")
		require.False(diag.HasErrors())

		diag = Configure(&c, f.Body, nil)
		require.False(diag.HasErrors())
		require.Equal(c.config.Name, "foo")
	})

	t.Run("invalid config", func(t *testing.T) {
		require := require.New(t)

		var c impl
		src := ``
		f, diag := hclparse.NewParser().ParseHCL([]byte(src), "test.hcl")
		require.False(diag.HasErrors())

		diag = Configure(&c, f.Body, nil)
		require.True(diag.HasErrors())
		require.Contains(diag.Error(), "is required")
	})

	t.Run("empty body", func(t *testing.T) {
		require := require.New(t)

		var s struct {
			Block struct {
				Label string   `hcl:",label"`
				Body  hcl.Body `hcl:",remain"`
			} `hcl:"block,block"`
		}

		src := `block "foo" {}`
		require.NoError(hclsimple.Decode("test.hcl", []byte(src), nil, &s))

		var c impl
		diag := Configure(&c, s.Block.Body, nil)
		require.True(diag.HasErrors())
		require.Contains(diag.Error(), "is required")
	})
}

func TestConfigure_nonImpl(t *testing.T) {
	t.Run("empty body", func(t *testing.T) {
		require := require.New(t)

		var s struct {
			Block struct {
				Label string   `hcl:",label"`
				Body  hcl.Body `hcl:",remain"`
			} `hcl:"block,block"`
		}

		src := `block "foo" {}`
		require.NoError(hclsimple.Decode("test.hcl", []byte(src), nil, &s))

		var c struct{}
		diag := Configure(&c, s.Block.Body, nil)
		require.False(diag.HasErrors())
	})

	t.Run("body", func(t *testing.T) {
		require := require.New(t)

		src := `name = "foo"`
		f, diag := hclparse.NewParser().ParseHCL([]byte(src), "test.hcl")
		require.False(diag.HasErrors())

		var c struct{}
		diag = Configure(&c, f.Body, nil)
		require.True(diag.HasErrors())
		t.Log(diag.Error())
	})
}

type testConfig struct {
	Name string `hcl:"name,attr"`
}

type impl struct{ config testConfig }

func (c *impl) Config() interface{} { return &c.config }
