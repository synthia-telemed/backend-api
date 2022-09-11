// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cache/client.go

// Package mock_cache_client is a generated GoMock package.
package mock_cache_client

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockClient) Get(ctx context.Context, key string, getAndDelete bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, getAndDelete)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(ctx, key, getAndDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), ctx, key, getAndDelete)
}

// HashGet mocks base method.
func (m *MockClient) HashGet(ctx context.Context, key, field string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashGet", ctx, key, field)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HashGet indicates an expected call of HashGet.
func (mr *MockClientMockRecorder) HashGet(ctx, key, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashGet", reflect.TypeOf((*MockClient)(nil).HashGet), ctx, key, field)
}

// HashSet mocks base method.
func (m *MockClient) HashSet(ctx context.Context, key string, kv map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashSet", ctx, key, kv)
	ret0, _ := ret[0].(error)
	return ret0
}

// HashSet indicates an expected call of HashSet.
func (mr *MockClientMockRecorder) HashSet(ctx, key, kv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashSet", reflect.TypeOf((*MockClient)(nil).HashSet), ctx, key, kv)
}

// MultipleGet mocks base method.
func (m *MockClient) MultipleGet(ctx context.Context, keys ...string) ([]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MultipleGet", varargs...)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultipleGet indicates an expected call of MultipleGet.
func (mr *MockClientMockRecorder) MultipleGet(ctx interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultipleGet", reflect.TypeOf((*MockClient)(nil).MultipleGet), varargs...)
}

// MultipleSet mocks base method.
func (m *MockClient) MultipleSet(ctx context.Context, kv map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultipleSet", ctx, kv)
	ret0, _ := ret[0].(error)
	return ret0
}

// MultipleSet indicates an expected call of MultipleSet.
func (mr *MockClientMockRecorder) MultipleSet(ctx, kv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultipleSet", reflect.TypeOf((*MockClient)(nil).MultipleSet), ctx, kv)
}

// Set mocks base method.
func (m *MockClient) Set(ctx context.Context, key, value string, expiredIn time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expiredIn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockClientMockRecorder) Set(ctx, key, value, expiredIn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockClient)(nil).Set), ctx, key, value, expiredIn)
}
