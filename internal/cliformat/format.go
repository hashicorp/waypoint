// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cliformat

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// FormatTableJson takes a waypoint-plugin-sdk.terminal Table and remaps it
// into a json output.
func FormatTableJson(t *terminal.Table) (string, error) {
	tableData := formatJsonMap(t)

	data, err := json.MarshalIndent(tableData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Takes a terminal Table and formats it into a map of key values to be used
// for formatting a JSON output response. It assumes any newline or comma separated
// values in the value entry are multiple items that should be formatted
// as a json array.
func formatJsonMap(t *terminal.Table) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, row := range t.Rows {
		c := map[string]interface{}{}

		for j, r := range row {
			// Remove any whitespacess in key
			header := strings.ReplaceAll(t.Headers[j], " ", "")
			// Lower case header key
			header = strings.ToLower(header)

			// This formatter assumes any strings that have a "\n" or a ", " in them
			// are multiple entries for a table. We join them back together so that
			// we can display them as a json array
			var vals []string
			if vals = strings.Split(r.Value, ", "); len(vals) > 1 {
				c[header] = vals
			} else if vals = strings.Split(r.Value, "\n"); len(vals) > 1 {
				c[header] = vals
			} else {
				c[header] = r.Value
			}
		}
		result = append(result, c)
	}

	return result
}
