package docker

// Image represents a Docker image.
type Image struct {
	// Image is the name of the image
	Image string

	// Tag is the tag associated with this image
	Tag string
}
