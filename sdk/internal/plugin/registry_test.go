package plugin

import (
	"testing"

	"github.com/mitchellh/devflow/sdk/component/mocks"
)

func TestRegistryConfig(t *testing.T) {
	mockV := &mockRegistryConfigurable{}
	testConfigurable(t, "registry", mockV, &mockV.Configurable)
}

type mockRegistryConfigurable struct {
	mocks.Registry
	mocks.Configurable
}
