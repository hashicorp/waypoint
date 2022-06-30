package installutil

import (
	"fmt"
	"github.com/distribution/distribution/v3/reference"
)

// DefaultODRImage returns the default Waypoint ODR image based on the
// supplied server image. We default the ODR image to the name of the server
// image with the `-odr` suffix attached to it.
func DefaultODRImage(serverImage string) (string, error) {
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
