// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint/internal/pkg/ctystructure"
)

// evalCtxTemplateProto adds template data to the eval context if tpl has a
// "TemplateData" field. This does nothing for any other type.
func evalCtxTemplateProto(ctx *hcl.EvalContext, key string, tpl proto.Message) error {
	// Get our template data field in our proto message. If we don't have one
	// then we don't bother doing anything.
	val := msgField(tpl, "TemplateData")
	if !val.IsValid() {
		return nil
	}

	valBytes, ok := val.Interface().([]byte)
	if !ok {
		return nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(valBytes, &raw); err != nil {
		return fmt.Errorf("Error decoding template data for %T: %w", tpl, err)
	}

	ctyVal, err := ctystructure.Object(raw)
	if err != nil {
		return err
	}

	if ctx.Variables == nil {
		ctx.Variables = map[string]cty.Value{}
	}
	ctx.Variables[key] = ctyVal

	return nil
}
