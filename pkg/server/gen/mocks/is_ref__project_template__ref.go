// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// isRef_ProjectTemplate_Ref is an autogenerated mock type for the isRef_ProjectTemplate_Ref type
type isRef_ProjectTemplate_Ref struct {
	mock.Mock
}

// isRef_ProjectTemplate_Ref provides a mock function with given fields:
func (_m *isRef_ProjectTemplate_Ref) isRef_ProjectTemplate_Ref() {
	_m.Called()
}

type mockConstructorTestingTnewIsRef_ProjectTemplate_Ref interface {
	mock.TestingT
	Cleanup(func())
}

// newIsRef_ProjectTemplate_Ref creates a new instance of isRef_ProjectTemplate_Ref. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newIsRef_ProjectTemplate_Ref(t mockConstructorTestingTnewIsRef_ProjectTemplate_Ref) *isRef_ProjectTemplate_Ref {
	mock := &isRef_ProjectTemplate_Ref{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
