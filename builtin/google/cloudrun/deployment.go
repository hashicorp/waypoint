package cloudrun

import (
	"fmt"

	"github.com/mitchellh/devflow/sdk/component"
)

// apiResource returns the GCP API "resource" string format for API calls.
func (d *Deployment) apiResource() string {
	return fmt.Sprintf("projects/%s/locations/%s/services/%s",
		d.Resource.Project,
		d.Resource.Location,
		d.Resource.Name,
	)
}

// apiName returns the GCP API "name" string format for API calls.
func (d *Deployment) apiName() string {
	return fmt.Sprintf("namespaces/%s/services/%s",
		d.Resource.Project,
		d.Resource.Name,
	)
}

var _ component.Deployment = (*Deployment)(nil)
