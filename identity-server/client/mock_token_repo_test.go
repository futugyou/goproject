// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/futugyousuzu/identity-server/client (interfaces: IOAuthTokenRepository)

// Package client_test is a generated GoMock package.
package client_test

import (
	context "context"
	reflect "reflect"

	client "github.com/futugyousuzu/identity-server/client"
	gomock "go.uber.org/mock/gomock"
)

// MockIOAuthTokenRepository is a mock of IOAuthTokenRepository interface.
type MockIOAuthTokenRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIOAuthTokenRepositoryMockRecorder
}

// MockIOAuthTokenRepositoryMockRecorder is the mock recorder for MockIOAuthTokenRepository.
type MockIOAuthTokenRepositoryMockRecorder struct {
	mock *MockIOAuthTokenRepository
}

// NewMockIOAuthTokenRepository creates a new mock instance.
func NewMockIOAuthTokenRepository(ctrl *gomock.Controller) *MockIOAuthTokenRepository {
	mock := &MockIOAuthTokenRepository{ctrl: ctrl}
	mock.recorder = &MockIOAuthTokenRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOAuthTokenRepository) EXPECT() *MockIOAuthTokenRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIOAuthTokenRepository) Get(arg0 context.Context, arg1 string) (*client.TokenModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*client.TokenModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIOAuthTokenRepositoryMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIOAuthTokenRepository)(nil).Get), arg0, arg1)
}

// Insert mocks base method.
func (m *MockIOAuthTokenRepository) Insert(arg0 context.Context, arg1 *client.TokenModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockIOAuthTokenRepositoryMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockIOAuthTokenRepository)(nil).Insert), arg0, arg1)
}

// Update mocks base method.
func (m *MockIOAuthTokenRepository) Update(arg0 context.Context, arg1 *client.TokenModel, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIOAuthTokenRepositoryMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIOAuthTokenRepository)(nil).Update), arg0, arg1, arg2)
}