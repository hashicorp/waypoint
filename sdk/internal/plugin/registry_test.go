package plugin

import (
	"testing"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/component/mocks"
)

func TestRegistryDynamicFunc_validateAuth(t *testing.T) {
	testDynamicFunc(t, "registry", &mockRegistryAuthenticator{}, func(v, f interface{}) {
		v.(*mockRegistryAuthenticator).Authenticator.On("ValidateAuthFunc").Return(f)
	}, func(raw interface{}) interface{} {
		return raw.(component.Authenticator).ValidateAuthFunc()
	})
}
func TestRegistryDynamicFunc_auth(t *testing.T) {
	testDynamicFunc(t, "registry", &mockRegistryAuthenticator{}, func(v, f interface{}) {
		v.(*mockRegistryAuthenticator).Authenticator.On("AuthFunc").Return(f)
	}, func(raw interface{}) interface{} {
		return raw.(component.Authenticator).AuthFunc()
	})
}

func TestRegistryConfig(t *testing.T) {
	mockV := &mockRegistryConfigurable{}
	testConfigurable(t, "registry", mockV, &mockV.Configurable)
}

type mockRegistryAuthenticator struct {
	mocks.Registry
	mocks.Authenticator
}

type mockRegistryConfigurable struct {
	mocks.Registry
	mocks.Configurable
}
