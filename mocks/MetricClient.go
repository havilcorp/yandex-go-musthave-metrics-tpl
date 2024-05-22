// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	metric_proto "github.com/havilcorp/yandex-go-musthave-metrics-tpl/pkg/proto/metric"
	mock "github.com/stretchr/testify/mock"
)

// MetricClient is an autogenerated mock type for the MetricClient type
type MetricClient struct {
	mock.Mock
}

// UpdateMetricBulk provides a mock function with given fields: ctx, in, opts
func (_m *MetricClient) UpdateMetricBulk(ctx context.Context, in *metric_proto.UpdateMetricBulkRequest, opts ...grpc.CallOption) (*metric_proto.UpdateMetricBulkResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMetricBulk")
	}

	var r0 *metric_proto.UpdateMetricBulkResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *metric_proto.UpdateMetricBulkRequest, ...grpc.CallOption) (*metric_proto.UpdateMetricBulkResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *metric_proto.UpdateMetricBulkRequest, ...grpc.CallOption) *metric_proto.UpdateMetricBulkResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metric_proto.UpdateMetricBulkResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *metric_proto.UpdateMetricBulkRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMetricClient creates a new instance of MetricClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetricClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetricClient {
	mock := &MetricClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}