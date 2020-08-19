package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
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
