package terraform

import (
	"encoding/json"
	"io"
)

// tfOutput is the type of an output from `terraform output -json`.
type tfOutput struct {
	Sensitive bool        `json:"sensitive"`
	Type      interface{} `json:"type"`
	Value     interface{} `json:"value"`
}

// parseOutputs parses the outputs from a reader that has the output from
// terraform output -json.
func parseOutputs(r io.Reader) (map[string]interface{}, error) {
	var outputs map[string]tfOutput
	if err := json.NewDecoder(r).Decode(&outputs); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for k, v := range outputs {
		result[k] = v.Value
	}

	return result, nil
}
