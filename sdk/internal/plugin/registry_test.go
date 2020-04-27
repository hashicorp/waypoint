package plugin

import (
	"testing"

	"github.com/hashicorp/waypoint/sdk/component/mocks"
)

func TestRegistryConfig(t *testing.T) {
	mockV := &mockRegistryConfigurable{}
	testConfigurable(t, "registry", mockV, &mockV.Configurable)
}

type mockRegistryConfigurable struct {
	mocks.Registry
	mocks.Configurable
}
