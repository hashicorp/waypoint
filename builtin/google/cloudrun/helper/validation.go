package helper

import (
	"fmt"
	"strings"
)

// ValidateImageName validates that that the specified image is in the gcr Docker Registry for this project
// Returns an error message when validation fails.
func ValidateImageName(image string, project string) error {
	// cloud run deployments must come from one of the following image registries
	var validRegistries = []string{
		"gcr.io",
		"us.gcr.io",
		"eu.gcr.io",
		"asia.gcr.io",
	}

	// check the image name has the valid parts
	parts := strings.Split(image, "/")
	if len(parts) != 3 {
		return fmt.Errorf("Invalid container image %s. Container images should be hosted in a Google Cloud registry for your project, i.e. gcr.io/%s/helloworld", image, project)
	}

	//check the registry is one which can be used with cloud run
	registryValid := false
	for _, r := range validRegistries {
		if r == parts[0] {
			registryValid = true
			break
		}
	}

	if !registryValid {
		return fmt.Errorf("Invalid container registry %s. Container images should be hosted in a valid Google Cloud registry e.g. %s", parts[0], strings.Join(parts, ","))
	}

	// check the project
	if parts[1] != project {
		return fmt.Errorf("Invalid container registry project %s. Container images should be hosted in Google Cloud registry for your project e.g. %s/%s/%s", parts[1], parts[0], project, parts[2])
	}

	return nil
}
