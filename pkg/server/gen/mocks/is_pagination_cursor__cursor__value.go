// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// isPaginationCursor_Cursor_Value is an autogenerated mock type for the isPaginationCursor_Cursor_Value type
type isPaginationCursor_Cursor_Value struct {
	mock.Mock
}

// isPaginationCursor_Cursor_Value provides a mock function with given fields:
func (_m *isPaginationCursor_Cursor_Value) isPaginationCursor_Cursor_Value() {
	_m.Called()
}

type mockConstructorTestingTnewIsPaginationCursor_Cursor_Value interface {
	mock.TestingT
	Cleanup(func())
}

// newIsPaginationCursor_Cursor_Value creates a new instance of isPaginationCursor_Cursor_Value. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newIsPaginationCursor_Cursor_Value(t mockConstructorTestingTnewIsPaginationCursor_Cursor_Value) *isPaginationCursor_Cursor_Value {
	mock := &isPaginationCursor_Cursor_Value{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
