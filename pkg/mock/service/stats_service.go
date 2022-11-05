// Code generated by MockGen. DO NOT EDIT.
// Source: ./stats_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	service "gosocialgraph/pkg/service"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStatsRunner is a mock of StatsRunner interface.
type MockStatsRunner struct {
	ctrl     *gomock.Controller
	recorder *MockStatsRunnerMockRecorder
}

// MockStatsRunnerMockRecorder is the mock recorder for MockStatsRunner.
type MockStatsRunnerMockRecorder struct {
	mock *MockStatsRunner
}

// NewMockStatsRunner creates a new mock instance.
func NewMockStatsRunner(ctrl *gomock.Controller) *MockStatsRunner {
	mock := &MockStatsRunner{ctrl: ctrl}
	mock.recorder = &MockStatsRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatsRunner) EXPECT() *MockStatsRunnerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockStatsRunner) Run(userID uuid.UUID) (service.StatsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", userID)
	ret0, _ := ret[0].(service.StatsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockStatsRunnerMockRecorder) Run(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockStatsRunner)(nil).Run), userID)
}
