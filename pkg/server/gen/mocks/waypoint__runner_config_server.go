// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gen "github.com/hashicorp/waypoint/pkg/server/gen"
	metadata "google.golang.org/grpc/metadata"

	mock "github.com/stretchr/testify/mock"
)

// Waypoint_RunnerConfigServer is an autogenerated mock type for the Waypoint_RunnerConfigServer type
type Waypoint_RunnerConfigServer struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *Waypoint_RunnerConfigServer) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// Recv provides a mock function with given fields:
func (_m *Waypoint_RunnerConfigServer) Recv() (*gen.RunnerConfigRequest, error) {
	ret := _m.Called()

	var r0 *gen.RunnerConfigRequest
	if rf, ok := ret.Get(0).(func() *gen.RunnerConfigRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.RunnerConfigRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RecvMsg provides a mock function with given fields: m
func (_m *Waypoint_RunnerConfigServer) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: _a0
func (_m *Waypoint_RunnerConfigServer) Send(_a0 *gen.RunnerConfigResponse) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gen.RunnerConfigResponse) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendHeader provides a mock function with given fields: _a0
func (_m *Waypoint_RunnerConfigServer) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMsg provides a mock function with given fields: m
func (_m *Waypoint_RunnerConfigServer) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetHeader provides a mock function with given fields: _a0
func (_m *Waypoint_RunnerConfigServer) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *Waypoint_RunnerConfigServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewWaypoint_RunnerConfigServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewWaypoint_RunnerConfigServer creates a new instance of Waypoint_RunnerConfigServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWaypoint_RunnerConfigServer(t mockConstructorTestingTNewWaypoint_RunnerConfigServer) *Waypoint_RunnerConfigServer {
	mock := &Waypoint_RunnerConfigServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
