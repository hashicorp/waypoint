package k8s

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"io/ioutil"
	"os"
	"testing"
)

// TestStartTask makes sure that we can use this Task Launcher Configuration - useful
// to test for example new fields and their parsing.

func TestStartTask(t *testing.T) {
	f, err := os.Open("./testdata/config.hcl")
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("unable to read file: %v", err)
	}

	hclFile, diags := hclsyntax.ParseConfig(fileContent, "plugin.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Fatalf("invalid HCL file: %v", diags.Error())
	}

	var hclCtx *hcl.EvalContext
	var tlc TaskLauncherConfig
	taskLauncher := TaskLauncher{
		config: tlc,
	}

	diag := component.Configure(&taskLauncher, hclFile.Body, hclCtx.NewChild())
	if diag.HasErrors() {
		t.Fatalf("diag failed: %v", diag.Error())
	}

	fmt.Printf("task_launcher_config=%+v\n", taskLauncher.config)
}
