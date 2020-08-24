package cloudrun

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	run "google.golang.org/api/run/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
		return fmt.Errorf("Invalid container image '%s'. Container images should be hosted in a Google Cloud registry for your project, i.e. 'gcr.io/%s/helloworld'", image, project)
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
		return fmt.Errorf("Invalid container registry '%s'. Container images should be hosted in a valid Google Cloud registry e.g. '%s'", parts[0], strings.Join(parts, ","))
	}

	// check the project
	if parts[1] != project {
		return fmt.Errorf("Invalid container registry project '%s'. Container images should be hosted in Google Cloud registry for your project e.g. '%s/%s/%s'", parts[1], parts[0], project, parts[2])
	}

	return nil
}

// ValidateRegionAvailable validates that the given GCP region is available for the project
func ValidateRegionAvailable(region string, gpcLocations []*run.Location) error {
	// keep a list of the regions so we can return a detailed error message
	regions := []string{}
	for _, r := range gpcLocations {
		if r.LocationId == region {
			return nil
		}

		regions = append(regions, r.LocationId)
	}

	return fmt.Errorf("The region '%s' is not available for this project, available regions are: '%s'", region, strings.Join(regions, ","))
}

var ErrInvalidMemoryValue = fmt.Errorf("Memory must be specified as a Kubernetes Quantity (https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/), and be less than '4Gi\n'")
var ErrInvalidCPUCount = fmt.Errorf("Invalid value for CPUCount, it is currently only possible to specify '1' or '2' CPUs\n")
var ErrInvalidRequestTimetout = fmt.Errorf("RequestTimeout must be greater than 0 and lets than 900\n")
var ErrInvalidMaxRequests = fmt.Errorf("MaxRequestsPerContainer must be greater than 0\n")
var ErrInvalidAutoscalingMax = fmt.Errorf("AutoScaling maximum must be larger than 0\n")

// ValidateConfig checks the deployment configuration for errors
func ValidateConfig(c Config) error {
	v := validator.New()
	v.RegisterValidation("kubernetes-memory", validateKubernetesMemory)

	err := v.Struct(c)

	if err != nil {
		errorMessage := ""
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Namespace() {
			case "Config.Capacity.Memory":
				errorMessage += ErrInvalidMemoryValue.Error()
			case "Config.Capacity.CPUCount":
				errorMessage += ErrInvalidCPUCount.Error()
			case "Config.Capacity.RequestTimeout":
				errorMessage += ErrInvalidRequestTimetout.Error()
			case "Config.Capacity.MaxRequestsPerContainer":
				errorMessage += ErrInvalidMaxRequests.Error()
			case "Config.AutoScaling.Max":
				errorMessage += ErrInvalidAutoscalingMax.Error()
			default:
				errorMessage += fmt.Sprintf("%s\n", err.Value())
			}
		}

		// if
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func validateKubernetesMemory(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	// if empty skip and use default value
	if v == "" {
		return true
	}

	// check that the value for memory is valid
	q, err := resource.ParseQuantity(v)
	if err != nil {
		return false
	}

	// check that quantity is max 4Gi
	var maxMemoryBytes int64 = 4294967296
	mi, _ := q.AsInt64()

	return mi <= maxMemoryBytes
}
