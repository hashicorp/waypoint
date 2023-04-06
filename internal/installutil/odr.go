// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package installutil

import (
	"fmt"

	"github.com/distribution/distribution/v3/reference"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// DeriveDefaultODRImage returns the default Waypoint ODR image based on the
// supplied server image. We default the ODR image to the name of the server
// image with the `-odr` suffix attached to it.
func DeriveDefaultODRImage(serverImage string) (string, error) {
	image, err := reference.Parse(serverImage)
	if err != nil {
		return "", fmt.Errorf("server image name %q is not a valid oci reference: %s", serverImage, err)
	}
	tagged, ok := image.(reference.Tagged)
	if !ok {
		return "", fmt.Errorf("server image doesn't have a tag specified. " +
			"Please specify a tag, for example `waypoint:latest`.")
	}

	tag := tagged.Tag()

	// Everything but the tag
	imageName := serverImage[0 : len(serverImage)-len(tag)-1]

	return fmt.Sprintf("%s-odr:%s", imageName, tag), nil
}

// NOTE: the server image is also used for static (non-ODR) runners.
// Static runners cannot use the ODR image.
const DefaultServerImage = "hashicorp/waypoint:latest"

// When we have a serverImage value to give to DeriveDefaultOdrImage,
// we should use that. When we don't, we can use this value
const DefaultODRImage = "hashicorp/waypoint-odr:latest"

func DefaultRunnerName(id string) string {
	return "waypoint-" + id + "-runner"
}

// An optional interface that the installer can implement to request
// an ondemand runner be registered.
type OnDemandRunnerConfigProvider interface {
	OnDemandRunnerConfig() *pb.OnDemandRunnerConfig
}
