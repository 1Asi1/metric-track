// Code generated by MockGen. DO NOT EDIT.
// Source: rpc/gen/metric_grpc.pb.go

// Package metricmock is a generated GoMock package.
package metricmock

import (
	context "context"
	reflect "reflect"

	gen "github.com/1Asi1/metric-track.git/rpc/gen"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockMetricGrpcClient is a mock of MetricGrpcClient interface.
type MockMetricGrpcClient struct {
	ctrl     *gomock.Controller
	recorder *MockMetricGrpcClientMockRecorder
}

// MockMetricGrpcClientMockRecorder is the mock recorder for MockMetricGrpcClient.
type MockMetricGrpcClientMockRecorder struct {
	mock *MockMetricGrpcClient
}

// NewMockMetricGrpcClient creates a new mock instance.
func NewMockMetricGrpcClient(ctrl *gomock.Controller) *MockMetricGrpcClient {
	mock := &MockMetricGrpcClient{ctrl: ctrl}
	mock.recorder = &MockMetricGrpcClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricGrpcClient) EXPECT() *MockMetricGrpcClientMockRecorder {
	return m.recorder
}

// Updates mocks base method.
func (m *MockMetricGrpcClient) Updates(ctx context.Context, in *gen.UpdatesRequest, opts ...grpc.CallOption) (*gen.UpdatesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Updates", varargs...)
	ret0, _ := ret[0].(*gen.UpdatesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Updates indicates an expected call of Updates.
func (mr *MockMetricGrpcClientMockRecorder) Updates(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Updates", reflect.TypeOf((*MockMetricGrpcClient)(nil).Updates), varargs...)
}

// MockMetricGrpcServer is a mock of MetricGrpcServer interface.
type MockMetricGrpcServer struct {
	ctrl     *gomock.Controller
	recorder *MockMetricGrpcServerMockRecorder
}

// MockMetricGrpcServerMockRecorder is the mock recorder for MockMetricGrpcServer.
type MockMetricGrpcServerMockRecorder struct {
	mock *MockMetricGrpcServer
}

// NewMockMetricGrpcServer creates a new mock instance.
func NewMockMetricGrpcServer(ctrl *gomock.Controller) *MockMetricGrpcServer {
	mock := &MockMetricGrpcServer{ctrl: ctrl}
	mock.recorder = &MockMetricGrpcServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricGrpcServer) EXPECT() *MockMetricGrpcServerMockRecorder {
	return m.recorder
}

// Updates mocks base method.
func (m *MockMetricGrpcServer) Updates(arg0 context.Context, arg1 *gen.UpdatesRequest) (*gen.UpdatesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Updates", arg0, arg1)
	ret0, _ := ret[0].(*gen.UpdatesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Updates indicates an expected call of Updates.
func (mr *MockMetricGrpcServerMockRecorder) Updates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Updates", reflect.TypeOf((*MockMetricGrpcServer)(nil).Updates), arg0, arg1)
}

// mustEmbedUnimplementedMetricGrpcServer mocks base method.
func (m *MockMetricGrpcServer) mustEmbedUnimplementedMetricGrpcServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedMetricGrpcServer")
}

// mustEmbedUnimplementedMetricGrpcServer indicates an expected call of mustEmbedUnimplementedMetricGrpcServer.
func (mr *MockMetricGrpcServerMockRecorder) mustEmbedUnimplementedMetricGrpcServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedMetricGrpcServer", reflect.TypeOf((*MockMetricGrpcServer)(nil).mustEmbedUnimplementedMetricGrpcServer))
}

// MockUnsafeMetricGrpcServer is a mock of UnsafeMetricGrpcServer interface.
type MockUnsafeMetricGrpcServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeMetricGrpcServerMockRecorder
}

// MockUnsafeMetricGrpcServerMockRecorder is the mock recorder for MockUnsafeMetricGrpcServer.
type MockUnsafeMetricGrpcServerMockRecorder struct {
	mock *MockUnsafeMetricGrpcServer
}

// NewMockUnsafeMetricGrpcServer creates a new mock instance.
func NewMockUnsafeMetricGrpcServer(ctrl *gomock.Controller) *MockUnsafeMetricGrpcServer {
	mock := &MockUnsafeMetricGrpcServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeMetricGrpcServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeMetricGrpcServer) EXPECT() *MockUnsafeMetricGrpcServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedMetricGrpcServer mocks base method.
func (m *MockUnsafeMetricGrpcServer) mustEmbedUnimplementedMetricGrpcServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedMetricGrpcServer")
}

// mustEmbedUnimplementedMetricGrpcServer indicates an expected call of mustEmbedUnimplementedMetricGrpcServer.
func (mr *MockUnsafeMetricGrpcServerMockRecorder) mustEmbedUnimplementedMetricGrpcServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedMetricGrpcServer", reflect.TypeOf((*MockUnsafeMetricGrpcServer)(nil).mustEmbedUnimplementedMetricGrpcServer))
}
