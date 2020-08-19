package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
	run "google.golang.org/api/run/v1"
)

func TestValidateImageReturnsErrorOnInvalidImageName(t *testing.T) {
	err := ValidateImageName("foo", "proj")
	require.Error(t, err)
}

func TestValidateImageReturnsErrorOnInvalidRegistry(t *testing.T) {
	err := ValidateImageName("foo/proj/image", "proj")
	require.Error(t, err)
}

func TestValidateImageReturnsErrorOnInvalidProject(t *testing.T) {
	err := ValidateImageName("gcr.io/proj2/image", "proj")
	require.Error(t, err)
}

func TestValidateImageReturnsNoErrorWhenValid(t *testing.T) {
	err := ValidateImageName("gcr.io/proj/image:latest", "proj")
	require.NoError(t, err)
}

var locations = []*run.Location{
	&run.Location{LocationId: "asia-east1"},
	&run.Location{LocationId: "asia-northeast1"},
}

func TestValidateRegionAvailableReturnsErrorWhenLocationNotAvailable(t *testing.T) {
	err := ValidateRegionAvailable("badlocation", locations)
	require.Error(t, err)
}

func TestValidateRegionAvailableReturnsNoErrorWhenLocationAvailable(t *testing.T) {
	err := ValidateRegionAvailable("asia-east1", locations)
	require.NoError(t, err)
}
