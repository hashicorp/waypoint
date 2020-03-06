package plugin

import (
	"testing"

	"github.com/mitchellh/devflow/sdk/component/mocks"
)

func TestPlatformConfig(t *testing.T) {
	mockV := &mockPlatformConfigurable{}
	testConfigurable(t, "platform", mockV, &mockV.Configurable)
}

type mockPlatformConfigurable struct {
	mocks.Platform
	mocks.Configurable
}
