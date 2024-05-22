// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DataBaseSaver is an autogenerated mock type for the DataBaseSaver type
type DataBaseSaver struct {
	mock.Mock
}

// Ping provides a mock function with given fields:
func (_m *DataBaseSaver) Ping() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ping")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDataBaseSaver creates a new instance of DataBaseSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataBaseSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataBaseSaver {
	mock := &DataBaseSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}