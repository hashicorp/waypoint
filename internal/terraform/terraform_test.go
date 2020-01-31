package terraform

import (
	"context"
	"os/exec"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/datadir"
)

//go:generate go-bindata -fs -nomemcopy -nometadata -pkg terraform -o bindata_test.go -prefix testdata/ testdata/...

// hasTerraform is used to determine if Terraform is installed. We only
// run certain tests if Terraform is available.
var hasTerraform bool

func init() {
	_, err := exec.LookPath("terraform")
	hasTerraform = err == nil
}

func TestTerraform(t *testing.T) {
	require := require.New(t)

	dir, closer := datadir.TestDir(t)
	defer closer()

	tf := &Terraform{
		Context:    context.Background(),
		Logger:     hclog.L(),
		Dir:        dir,
		ConfigFS:   AssetFile(),
		ConfigPath: "basic",
		Vars: map[string]interface{}{
			"number": 12,
		},
	}

	outputs, err := tf.Apply()
	require.NoError(err)
	require.Equal(outputs, map[string]interface{}{
		"double": float64(24),
	})
}
