// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package manifest

import (
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	sigyaml "sigs.k8s.io/yaml"
)

// A manifest represents a multi-resource Kubernetes resource description
// such as from YAML.
type Manifest struct {
	Resources []*Resource
}

// Resource is a generic resource. This is meant to be able to match
// any Kubernetes resource. Callers are expected to use the metadata to
// further query the object if more information is required.
type Resource struct {
	APIVersion string `mapstructure:"apiVersion"`
	Kind       string `mapstructure:"kind"`
	Metadata   struct {
		Name        string
		Namespace   string
		Labels      map[string]string
		Annotations map[string]string
	} `mapstructure:"metadata"`

	// RawYAML is the raw document YAML for this resource.
	RawYAML []byte

	// RawJSON is the RawYAML converted to JSON using the k8s.io YAML
	// library. We do this because all the API structures are tagged with
	// a "json" tag so this lets you unmarshal into a richer official
	// Kubernetes structure directly.
	RawJSON []byte
}

// FullKind returns the Kind of this resource with the APIVersion as a prefix.
// This will not prefix core types (such as Secret).
func (r *Resource) FullKind() string {
	if r.APIVersion == "" {
		return r.Kind
	}

	return fmt.Sprintf("%s/%s", r.APIVersion, r.Kind)
}

// Parse parses multi-document YAML contents into a manifest.
func Parse(r io.Reader) (*Manifest, error) {
	dec := yaml.NewDecoder(r)

	var m Manifest
	for {
		// Decode the raw document
		var raw map[string]interface{}
		err := dec.Decode(&raw)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Try to form a richer resource from this.
		var r Resource
		if err := mapstructure.Decode(raw, &r); err != nil {
			return nil, err
		}

		// Remarshal as YAML so we can get the single resource.
		r.RawYAML, err = yaml.Marshal(raw)
		if err != nil {
			return nil, err
		}

		// Remarshal as JSON. See the struct docs on why we have both.
		r.RawJSON, err = sigyaml.YAMLToJSON(r.RawYAML)
		if err != nil {
			return nil, err
		}

		m.Resources = append(m.Resources, &r)
	}

	return &m, nil
}
