// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// isListWorkspacesRequest_Scope is an autogenerated mock type for the isListWorkspacesRequest_Scope type
type isListWorkspacesRequest_Scope struct {
	mock.Mock
}

// isListWorkspacesRequest_Scope provides a mock function with given fields:
func (_m *isListWorkspacesRequest_Scope) isListWorkspacesRequest_Scope() {
	_m.Called()
}

type mockConstructorTestingTnewIsListWorkspacesRequest_Scope interface {
	mock.TestingT
	Cleanup(func())
}

// newIsListWorkspacesRequest_Scope creates a new instance of isListWorkspacesRequest_Scope. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newIsListWorkspacesRequest_Scope(t mockConstructorTestingTnewIsListWorkspacesRequest_Scope) *isListWorkspacesRequest_Scope {
	mock := &isListWorkspacesRequest_Scope{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
