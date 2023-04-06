// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package k8sstatus

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/builtin/k8s/internal/manifest"
)

func FromManifest(
	ctx context.Context,
	cs *kubernetes.Clientset,
	m *manifest.Manifest,
) (*sdk.StatusReport, error) {
	var report sdk.StatusReport

	// Go through each resource. All resources in the manifest will be present
	// in the final report even if we have no idea what they are (i.e. CRDs).
	for _, r := range m.Resources {
		var resource *sdk.StatusReport_Resource
		var err error

		// Unknown type
		resource, err = unknownResource(r)

		// If there was an error, return the error.
		if err != nil {
			return nil, err
		}

		// Add this resource to our report
		report.Resources = append(report.Resources, resource)
	}

	// Set additional metadata on the report
	report.GeneratedTime = timestamppb.Now()
	report.External = true

	return &report, nil
}

// unknownReource takes an unknown manifest resource and turns it into a
// status report resource.
func unknownResource(resource *manifest.Resource) (*sdk.StatusReport_Resource, error) {
	var r sdk.StatusReport_Resource
	r.Id = resource.Metadata.Name
	r.Platform = "kubernetes"
	r.Name = resource.Metadata.Name
	r.Type = resource.Kind
	r.StateJson = string(resource.RawJSON)
	return &r, nil
}
