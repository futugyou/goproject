// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/futugyousuzu/identity-server/token (interfaces: IJwksService)

// Package token_test is a generated GoMock package.
package token_test

import (
	context "context"
	reflect "reflect"

	jwk "github.com/lestrrat-go/jwx/v2/jwk"
	gomock "go.uber.org/mock/gomock"
)

// MockIJwksService is a mock of IJwksService interface.
type MockIJwksService struct {
	ctrl     *gomock.Controller
	recorder *MockIJwksServiceMockRecorder
}

// MockIJwksServiceMockRecorder is the mock recorder for MockIJwksService.
type MockIJwksServiceMockRecorder struct {
	mock *MockIJwksService
}

// NewMockIJwksService creates a new mock instance.
func NewMockIJwksService(ctrl *gomock.Controller) *MockIJwksService {
	mock := &MockIJwksService{ctrl: ctrl}
	mock.recorder = &MockIJwksServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIJwksService) EXPECT() *MockIJwksServiceMockRecorder {
	return m.recorder
}

// CreateJwks mocks base method.
func (m *MockIJwksService) CreateJwks(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJwks", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateJwks indicates an expected call of CreateJwks.
func (mr *MockIJwksServiceMockRecorder) CreateJwks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJwks", reflect.TypeOf((*MockIJwksService)(nil).CreateJwks), arg0, arg1)
}

// GetJwkByKeyID mocks base method.
func (m *MockIJwksService) GetJwkByKeyID(arg0 context.Context, arg1 string) (jwk.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJwkByKeyID", arg0, arg1)
	ret0, _ := ret[0].(jwk.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJwkByKeyID indicates an expected call of GetJwkByKeyID.
func (mr *MockIJwksServiceMockRecorder) GetJwkByKeyID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJwkByKeyID", reflect.TypeOf((*MockIJwksService)(nil).GetJwkByKeyID), arg0, arg1)
}

// GetPublicJwks mocks base method.
func (m *MockIJwksService) GetPublicJwks(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicJwks", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicJwks indicates an expected call of GetPublicJwks.
func (mr *MockIJwksServiceMockRecorder) GetPublicJwks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicJwks", reflect.TypeOf((*MockIJwksService)(nil).GetPublicJwks), arg0)
}
