// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"

	mock "github.com/stretchr/testify/mock"
)

// MetricSaver is an autogenerated mock type for the MetricSaver type
type MetricSaver struct {
	mock.Mock
}

// AddCounterBulk provides a mock function with given fields: ctx, list
func (_m *MetricSaver) AddCounterBulk(ctx context.Context, list []domain.Counter) error {
	ret := _m.Called(ctx, list)

	if len(ret) == 0 {
		panic("no return value specified for AddCounterBulk")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.Counter) error); ok {
		r0 = rf(ctx, list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddGaugeBulk provides a mock function with given fields: ctx, list
func (_m *MetricSaver) AddGaugeBulk(ctx context.Context, list []domain.Gauge) error {
	ret := _m.Called(ctx, list)

	if len(ret) == 0 {
		panic("no return value specified for AddGaugeBulk")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.Gauge) error); ok {
		r0 = rf(ctx, list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMetricSaver creates a new instance of MetricSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetricSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetricSaver {
	mock := &MetricSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
