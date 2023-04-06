// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aci

import (
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/go-playground/validator"
)

func validateLocationAvailable(location string, validLocations []string) error {
	for _, loc := range validLocations {
		if loc == location {
			return nil
		}
	}

	return fmt.Errorf("Location %s, is not a valid location for this subscription. Valid locations are: %s", location, validLocations)
}

var errInvalidMemoryValue = fmt.Errorf("Memory allocated to a Cloud run instance must a minimum of 512MB and less than 16384MB (16GB)\n")
var errInvalidCPUCount = fmt.Errorf("Invalid value for CPUCount, it is currently only possible to specify '1-4' CPUs\n")
var errInvalidVolume = fmt.Errorf("Container instance volumes must have one of 'azure_file_share' or 'git_repo' fields set\n")

func validateConfig(c Config) error {
	v := validator.New()
	v.RegisterStructValidation(validationVolumeStruct, Volume{})

	err := v.Struct(c)

	if err != nil {
		errorMessage := ""
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Namespace() {
			case "Config.Capacity.Memory":
				errorMessage += errInvalidMemoryValue.Error()
			case "Config.Capacity.CPUCount":
				errorMessage += errInvalidCPUCount.Error()
			}

			if err.Tag() == "one_of" {
				errorMessage += errInvalidVolume.Error()
			}
		}

		return fmt.Errorf(errorMessage)
	}

	return nil
}

func validationVolumeStruct(sl validator.StructLevel) {
	vol := sl.Current().Interface().(Volume)

	// Error if both set
	if vol.AzureFileShare != nil && vol.GitRepoVolume != nil {
		sl.ReportError(vol, "volume", "Volume", "one_of", "")
	}

	// Error if none
	if vol.AzureFileShare == nil && vol.GitRepoVolume == nil {
		sl.ReportError(vol, "volume", "Volume", "one_of", "")
	}
}

// return the server component from an image name
func parseDockerServer(image string) string {
	n, err := reference.ParseNamed(image)
	if err != nil {
		return "docker.io"
	}

	d := reference.Domain(n)
	if d == "" {
		// no domain, convention is main docker repo
		d = "docker.io"
	}

	return d
}
