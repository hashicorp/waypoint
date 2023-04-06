// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package formatter

import (
	"fmt"
	"sort"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var (
	// Variable value sources
	// listed in descending precedence order for ease of reference
	SourceCLI     = "cli"
	SourceFile    = "file"
	SourceEnv     = "env"
	SourceVCS     = "vcs"
	SourceServer  = "server"
	SourceDynamic = "dynamic"
	SourceDefault = "default"
	SourceUnknown = "unknown"

	fromFVtoSource = map[pb.Variable_FinalValue_Source]string{
		pb.Variable_FinalValue_UNKNOWN: SourceUnknown,
		pb.Variable_FinalValue_DEFAULT: SourceDefault,
		pb.Variable_FinalValue_FILE:    SourceFile,
		pb.Variable_FinalValue_CLI:     SourceCLI,
		pb.Variable_FinalValue_ENV:     SourceEnv,
		pb.Variable_FinalValue_VCS:     SourceVCS,
		pb.Variable_FinalValue_SERVER:  SourceServer,
		pb.Variable_FinalValue_DYNAMIC: SourceDynamic,
	}
)

type Output struct {
	Value  string
	Type   string
	Source string
}

func ValuesForOutput(values map[string]*pb.Variable_FinalValue) map[string]*Output {
	outputs := make(map[string]*Output, len(values))
	// sort alphabetically for joy
	inputVars := make([]string, 0, len(values))
	for iv := range values {
		inputVars = append(inputVars, iv)
	}
	sort.Strings(inputVars)
	for _, iv := range inputVars {
		outputs[iv] = &Output{}
		if o, ok := outputs[iv]; ok {
			// move value and inferred type into strings for outputting
			switch vt := values[iv].Value.(type) {
			case *pb.Variable_FinalValue_Sensitive:
				o.Value = vt.Sensitive
				o.Type = "sensitive"
			case *pb.Variable_FinalValue_Str:
				o.Value = vt.Str
				o.Type = "string"
			case *pb.Variable_FinalValue_Bool:
				o.Value = fmt.Sprintf("%t", vt.Bool)
				o.Type = "bool"
			case *pb.Variable_FinalValue_Num:
				o.Value = fmt.Sprintf("%d", vt.Num)
				o.Type = "int"
			case *pb.Variable_FinalValue_Hcl:
				o.Value = vt.Hcl
				o.Type = "complex"
			}

			o.Source = fromFVtoSource[values[iv].Source]

		}
	}

	return outputs
}
