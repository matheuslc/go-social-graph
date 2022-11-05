// Code generated by MockGen. DO NOT EDIT.
// Source: ./user_post_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	service "gosocialgraph/pkg/service"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserPostRunner is a mock of UserPostRunner interface.
type MockUserPostRunner struct {
	ctrl     *gomock.Controller
	recorder *MockUserPostRunnerMockRecorder
}

// MockUserPostRunnerMockRecorder is the mock recorder for MockUserPostRunner.
type MockUserPostRunnerMockRecorder struct {
	mock *MockUserPostRunner
}

// NewMockUserPostRunner creates a new mock instance.
func NewMockUserPostRunner(ctrl *gomock.Controller) *MockUserPostRunner {
	mock := &MockUserPostRunner{ctrl: ctrl}
	mock.recorder = &MockUserPostRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserPostRunner) EXPECT() *MockUserPostRunnerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockUserPostRunner) Run(intent service.UserPostsIntent) (service.UserPostResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", intent)
	ret0, _ := ret[0].(service.UserPostResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockUserPostRunnerMockRecorder) Run(intent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockUserPostRunner)(nil).Run), intent)
}
