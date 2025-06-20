// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package kayakv1connect

import (
	"context"

	"connectrpc.com/connect"
	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	mock "github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NewMockKayakServiceClient creates a new instance of MockKayakServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKayakServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKayakServiceClient {
	mock := &MockKayakServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockKayakServiceClient is an autogenerated mock type for the KayakServiceClient type
type MockKayakServiceClient struct {
	mock.Mock
}

type MockKayakServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKayakServiceClient) EXPECT() *MockKayakServiceClient_Expecter {
	return &MockKayakServiceClient_Expecter{mock: &_m.Mock}
}

// Apply provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) Apply(context1 context.Context, request *connect.Request[kayakv1.ApplyRequest]) (*connect.Response[kayakv1.ApplyResponse], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for Apply")
	}

	var r0 *connect.Response[kayakv1.ApplyResponse]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.ApplyRequest]) (*connect.Response[kayakv1.ApplyResponse], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.ApplyRequest]) *connect.Response[kayakv1.ApplyResponse]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[kayakv1.ApplyResponse])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.ApplyRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_Apply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Apply'
type MockKayakServiceClient_Apply_Call struct {
	*mock.Call
}

// Apply is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) Apply(context1 interface{}, request interface{}) *MockKayakServiceClient_Apply_Call {
	return &MockKayakServiceClient_Apply_Call{Call: _e.mock.On("Apply", context1, request)}
}

func (_c *MockKayakServiceClient_Apply_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.ApplyRequest])) *MockKayakServiceClient_Apply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.ApplyRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_Apply_Call) Return(response *connect.Response[kayakv1.ApplyResponse], err error) *MockKayakServiceClient_Apply_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_Apply_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.ApplyRequest]) (*connect.Response[kayakv1.ApplyResponse], error)) *MockKayakServiceClient_Apply_Call {
	_c.Call.Return(run)
	return _c
}

// CommitRecord provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) CommitRecord(context1 context.Context, request *connect.Request[kayakv1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for CommitRecord")
	}

	var r0 *connect.Response[emptypb.Empty]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.CommitRecordRequest]) *connect.Response[emptypb.Empty]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[emptypb.Empty])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.CommitRecordRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_CommitRecord_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommitRecord'
type MockKayakServiceClient_CommitRecord_Call struct {
	*mock.Call
}

// CommitRecord is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) CommitRecord(context1 interface{}, request interface{}) *MockKayakServiceClient_CommitRecord_Call {
	return &MockKayakServiceClient_CommitRecord_Call{Call: _e.mock.On("CommitRecord", context1, request)}
}

func (_c *MockKayakServiceClient_CommitRecord_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.CommitRecordRequest])) *MockKayakServiceClient_CommitRecord_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.CommitRecordRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_CommitRecord_Call) Return(response *connect.Response[emptypb.Empty], err error) *MockKayakServiceClient_CommitRecord_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_CommitRecord_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error)) *MockKayakServiceClient_CommitRecord_Call {
	_c.Call.Return(run)
	return _c
}

// CreateStream provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) CreateStream(context1 context.Context, request *connect.Request[kayakv1.CreateStreamRequest]) (*connect.Response[emptypb.Empty], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for CreateStream")
	}

	var r0 *connect.Response[emptypb.Empty]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.CreateStreamRequest]) (*connect.Response[emptypb.Empty], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.CreateStreamRequest]) *connect.Response[emptypb.Empty]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[emptypb.Empty])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.CreateStreamRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_CreateStream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateStream'
type MockKayakServiceClient_CreateStream_Call struct {
	*mock.Call
}

// CreateStream is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) CreateStream(context1 interface{}, request interface{}) *MockKayakServiceClient_CreateStream_Call {
	return &MockKayakServiceClient_CreateStream_Call{Call: _e.mock.On("CreateStream", context1, request)}
}

func (_c *MockKayakServiceClient_CreateStream_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.CreateStreamRequest])) *MockKayakServiceClient_CreateStream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.CreateStreamRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_CreateStream_Call) Return(response *connect.Response[emptypb.Empty], err error) *MockKayakServiceClient_CreateStream_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_CreateStream_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.CreateStreamRequest]) (*connect.Response[emptypb.Empty], error)) *MockKayakServiceClient_CreateStream_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteStream provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) DeleteStream(context1 context.Context, request *connect.Request[kayakv1.DeleteStreamRequest]) (*connect.Response[emptypb.Empty], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for DeleteStream")
	}

	var r0 *connect.Response[emptypb.Empty]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.DeleteStreamRequest]) (*connect.Response[emptypb.Empty], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.DeleteStreamRequest]) *connect.Response[emptypb.Empty]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[emptypb.Empty])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.DeleteStreamRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_DeleteStream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteStream'
type MockKayakServiceClient_DeleteStream_Call struct {
	*mock.Call
}

// DeleteStream is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) DeleteStream(context1 interface{}, request interface{}) *MockKayakServiceClient_DeleteStream_Call {
	return &MockKayakServiceClient_DeleteStream_Call{Call: _e.mock.On("DeleteStream", context1, request)}
}

func (_c *MockKayakServiceClient_DeleteStream_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.DeleteStreamRequest])) *MockKayakServiceClient_DeleteStream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.DeleteStreamRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_DeleteStream_Call) Return(response *connect.Response[emptypb.Empty], err error) *MockKayakServiceClient_DeleteStream_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_DeleteStream_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.DeleteStreamRequest]) (*connect.Response[emptypb.Empty], error)) *MockKayakServiceClient_DeleteStream_Call {
	_c.Call.Return(run)
	return _c
}

// DeregisterWorker provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) DeregisterWorker(context1 context.Context, request *connect.Request[kayakv1.DeregisterWorkerRequest]) (*connect.Response[emptypb.Empty], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for DeregisterWorker")
	}

	var r0 *connect.Response[emptypb.Empty]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.DeregisterWorkerRequest]) (*connect.Response[emptypb.Empty], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.DeregisterWorkerRequest]) *connect.Response[emptypb.Empty]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[emptypb.Empty])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.DeregisterWorkerRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_DeregisterWorker_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeregisterWorker'
type MockKayakServiceClient_DeregisterWorker_Call struct {
	*mock.Call
}

// DeregisterWorker is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) DeregisterWorker(context1 interface{}, request interface{}) *MockKayakServiceClient_DeregisterWorker_Call {
	return &MockKayakServiceClient_DeregisterWorker_Call{Call: _e.mock.On("DeregisterWorker", context1, request)}
}

func (_c *MockKayakServiceClient_DeregisterWorker_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.DeregisterWorkerRequest])) *MockKayakServiceClient_DeregisterWorker_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.DeregisterWorkerRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_DeregisterWorker_Call) Return(response *connect.Response[emptypb.Empty], err error) *MockKayakServiceClient_DeregisterWorker_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_DeregisterWorker_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.DeregisterWorkerRequest]) (*connect.Response[emptypb.Empty], error)) *MockKayakServiceClient_DeregisterWorker_Call {
	_c.Call.Return(run)
	return _c
}

// FetchRecords provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) FetchRecords(context1 context.Context, request *connect.Request[kayakv1.FetchRecordsRequest]) (*connect.Response[kayakv1.FetchRecordsResponse], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for FetchRecords")
	}

	var r0 *connect.Response[kayakv1.FetchRecordsResponse]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.FetchRecordsRequest]) (*connect.Response[kayakv1.FetchRecordsResponse], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.FetchRecordsRequest]) *connect.Response[kayakv1.FetchRecordsResponse]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[kayakv1.FetchRecordsResponse])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.FetchRecordsRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_FetchRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchRecords'
type MockKayakServiceClient_FetchRecords_Call struct {
	*mock.Call
}

// FetchRecords is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) FetchRecords(context1 interface{}, request interface{}) *MockKayakServiceClient_FetchRecords_Call {
	return &MockKayakServiceClient_FetchRecords_Call{Call: _e.mock.On("FetchRecords", context1, request)}
}

func (_c *MockKayakServiceClient_FetchRecords_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.FetchRecordsRequest])) *MockKayakServiceClient_FetchRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.FetchRecordsRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_FetchRecords_Call) Return(response *connect.Response[kayakv1.FetchRecordsResponse], err error) *MockKayakServiceClient_FetchRecords_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_FetchRecords_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.FetchRecordsRequest]) (*connect.Response[kayakv1.FetchRecordsResponse], error)) *MockKayakServiceClient_FetchRecords_Call {
	_c.Call.Return(run)
	return _c
}

// GetRecords provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) GetRecords(context1 context.Context, request *connect.Request[kayakv1.GetRecordsRequest]) (*connect.Response[kayakv1.GetRecordsResponse], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for GetRecords")
	}

	var r0 *connect.Response[kayakv1.GetRecordsResponse]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.GetRecordsRequest]) (*connect.Response[kayakv1.GetRecordsResponse], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.GetRecordsRequest]) *connect.Response[kayakv1.GetRecordsResponse]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[kayakv1.GetRecordsResponse])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.GetRecordsRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_GetRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRecords'
type MockKayakServiceClient_GetRecords_Call struct {
	*mock.Call
}

// GetRecords is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) GetRecords(context1 interface{}, request interface{}) *MockKayakServiceClient_GetRecords_Call {
	return &MockKayakServiceClient_GetRecords_Call{Call: _e.mock.On("GetRecords", context1, request)}
}

func (_c *MockKayakServiceClient_GetRecords_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.GetRecordsRequest])) *MockKayakServiceClient_GetRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.GetRecordsRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_GetRecords_Call) Return(response *connect.Response[kayakv1.GetRecordsResponse], err error) *MockKayakServiceClient_GetRecords_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_GetRecords_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.GetRecordsRequest]) (*connect.Response[kayakv1.GetRecordsResponse], error)) *MockKayakServiceClient_GetRecords_Call {
	_c.Call.Return(run)
	return _c
}

// GetStream provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) GetStream(context1 context.Context, request *connect.Request[kayakv1.GetStreamRequest]) (*connect.Response[kayakv1.GetStreamResponse], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for GetStream")
	}

	var r0 *connect.Response[kayakv1.GetStreamResponse]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.GetStreamRequest]) (*connect.Response[kayakv1.GetStreamResponse], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.GetStreamRequest]) *connect.Response[kayakv1.GetStreamResponse]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[kayakv1.GetStreamResponse])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.GetStreamRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_GetStream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStream'
type MockKayakServiceClient_GetStream_Call struct {
	*mock.Call
}

// GetStream is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) GetStream(context1 interface{}, request interface{}) *MockKayakServiceClient_GetStream_Call {
	return &MockKayakServiceClient_GetStream_Call{Call: _e.mock.On("GetStream", context1, request)}
}

func (_c *MockKayakServiceClient_GetStream_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.GetStreamRequest])) *MockKayakServiceClient_GetStream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.GetStreamRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_GetStream_Call) Return(response *connect.Response[kayakv1.GetStreamResponse], err error) *MockKayakServiceClient_GetStream_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_GetStream_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.GetStreamRequest]) (*connect.Response[kayakv1.GetStreamResponse], error)) *MockKayakServiceClient_GetStream_Call {
	_c.Call.Return(run)
	return _c
}

// GetStreams provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) GetStreams(context1 context.Context, request *connect.Request[kayakv1.GetStreamsRequest]) (*connect.Response[kayakv1.GetStreamsResponse], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for GetStreams")
	}

	var r0 *connect.Response[kayakv1.GetStreamsResponse]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.GetStreamsRequest]) (*connect.Response[kayakv1.GetStreamsResponse], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.GetStreamsRequest]) *connect.Response[kayakv1.GetStreamsResponse]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[kayakv1.GetStreamsResponse])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.GetStreamsRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_GetStreams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStreams'
type MockKayakServiceClient_GetStreams_Call struct {
	*mock.Call
}

// GetStreams is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) GetStreams(context1 interface{}, request interface{}) *MockKayakServiceClient_GetStreams_Call {
	return &MockKayakServiceClient_GetStreams_Call{Call: _e.mock.On("GetStreams", context1, request)}
}

func (_c *MockKayakServiceClient_GetStreams_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.GetStreamsRequest])) *MockKayakServiceClient_GetStreams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.GetStreamsRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_GetStreams_Call) Return(response *connect.Response[kayakv1.GetStreamsResponse], err error) *MockKayakServiceClient_GetStreams_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_GetStreams_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.GetStreamsRequest]) (*connect.Response[kayakv1.GetStreamsResponse], error)) *MockKayakServiceClient_GetStreams_Call {
	_c.Call.Return(run)
	return _c
}

// PutRecords provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) PutRecords(context1 context.Context, request *connect.Request[kayakv1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for PutRecords")
	}

	var r0 *connect.Response[emptypb.Empty]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.PutRecordsRequest]) *connect.Response[emptypb.Empty]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[emptypb.Empty])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.PutRecordsRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_PutRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutRecords'
type MockKayakServiceClient_PutRecords_Call struct {
	*mock.Call
}

// PutRecords is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) PutRecords(context1 interface{}, request interface{}) *MockKayakServiceClient_PutRecords_Call {
	return &MockKayakServiceClient_PutRecords_Call{Call: _e.mock.On("PutRecords", context1, request)}
}

func (_c *MockKayakServiceClient_PutRecords_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.PutRecordsRequest])) *MockKayakServiceClient_PutRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.PutRecordsRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_PutRecords_Call) Return(response *connect.Response[emptypb.Empty], err error) *MockKayakServiceClient_PutRecords_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_PutRecords_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error)) *MockKayakServiceClient_PutRecords_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterWorker provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) RegisterWorker(context1 context.Context, request *connect.Request[kayakv1.RegisterWorkerRequest]) (*connect.Response[kayakv1.RegisterWorkerResponse], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for RegisterWorker")
	}

	var r0 *connect.Response[kayakv1.RegisterWorkerResponse]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.RegisterWorkerRequest]) (*connect.Response[kayakv1.RegisterWorkerResponse], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.RegisterWorkerRequest]) *connect.Response[kayakv1.RegisterWorkerResponse]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[kayakv1.RegisterWorkerResponse])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.RegisterWorkerRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_RegisterWorker_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterWorker'
type MockKayakServiceClient_RegisterWorker_Call struct {
	*mock.Call
}

// RegisterWorker is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) RegisterWorker(context1 interface{}, request interface{}) *MockKayakServiceClient_RegisterWorker_Call {
	return &MockKayakServiceClient_RegisterWorker_Call{Call: _e.mock.On("RegisterWorker", context1, request)}
}

func (_c *MockKayakServiceClient_RegisterWorker_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.RegisterWorkerRequest])) *MockKayakServiceClient_RegisterWorker_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.RegisterWorkerRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_RegisterWorker_Call) Return(response *connect.Response[kayakv1.RegisterWorkerResponse], err error) *MockKayakServiceClient_RegisterWorker_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_RegisterWorker_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.RegisterWorkerRequest]) (*connect.Response[kayakv1.RegisterWorkerResponse], error)) *MockKayakServiceClient_RegisterWorker_Call {
	_c.Call.Return(run)
	return _c
}

// RenewRegistration provides a mock function for the type MockKayakServiceClient
func (_mock *MockKayakServiceClient) RenewRegistration(context1 context.Context, request *connect.Request[kayakv1.RenewRegistrationRequest]) (*connect.Response[emptypb.Empty], error) {
	ret := _mock.Called(context1, request)

	if len(ret) == 0 {
		panic("no return value specified for RenewRegistration")
	}

	var r0 *connect.Response[emptypb.Empty]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.RenewRegistrationRequest]) (*connect.Response[emptypb.Empty], error)); ok {
		return returnFunc(context1, request)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *connect.Request[kayakv1.RenewRegistrationRequest]) *connect.Response[emptypb.Empty]); ok {
		r0 = returnFunc(context1, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*connect.Response[emptypb.Empty])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *connect.Request[kayakv1.RenewRegistrationRequest]) error); ok {
		r1 = returnFunc(context1, request)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockKayakServiceClient_RenewRegistration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenewRegistration'
type MockKayakServiceClient_RenewRegistration_Call struct {
	*mock.Call
}

// RenewRegistration is a helper method to define mock.On call
//   - context1
//   - request
func (_e *MockKayakServiceClient_Expecter) RenewRegistration(context1 interface{}, request interface{}) *MockKayakServiceClient_RenewRegistration_Call {
	return &MockKayakServiceClient_RenewRegistration_Call{Call: _e.mock.On("RenewRegistration", context1, request)}
}

func (_c *MockKayakServiceClient_RenewRegistration_Call) Run(run func(context1 context.Context, request *connect.Request[kayakv1.RenewRegistrationRequest])) *MockKayakServiceClient_RenewRegistration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*connect.Request[kayakv1.RenewRegistrationRequest]))
	})
	return _c
}

func (_c *MockKayakServiceClient_RenewRegistration_Call) Return(response *connect.Response[emptypb.Empty], err error) *MockKayakServiceClient_RenewRegistration_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockKayakServiceClient_RenewRegistration_Call) RunAndReturn(run func(context1 context.Context, request *connect.Request[kayakv1.RenewRegistrationRequest]) (*connect.Response[emptypb.Empty], error)) *MockKayakServiceClient_RenewRegistration_Call {
	_c.Call.Return(run)
	return _c
}
