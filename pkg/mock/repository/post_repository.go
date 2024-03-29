// Code generated by MockGen. DO NOT EDIT.
// Source: ./post_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "gosocialgraph/pkg/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPostWriter is a mock of PostWriter interface.
type MockPostWriter struct {
	ctrl     *gomock.Controller
	recorder *MockPostWriterMockRecorder
}

// MockPostWriterMockRecorder is the mock recorder for MockPostWriter.
type MockPostWriterMockRecorder struct {
	mock *MockPostWriter
}

// NewMockPostWriter creates a new mock instance.
func NewMockPostWriter(ctrl *gomock.Controller) *MockPostWriter {
	mock := &MockPostWriter{ctrl: ctrl}
	mock.recorder = &MockPostWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostWriter) EXPECT() *MockPostWriterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPostWriter) Create(ctx context.Context, userID, content string) (entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userID, content)
	ret0, _ := ret[0].(entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPostWriterMockRecorder) Create(ctx, userID, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPostWriter)(nil).Create), ctx, userID, content)
}

// MockReposter is a mock of Reposter interface.
type MockReposter struct {
	ctrl     *gomock.Controller
	recorder *MockReposterMockRecorder
}

// MockReposterMockRecorder is the mock recorder for MockReposter.
type MockReposterMockRecorder struct {
	mock *MockReposter
}

// NewMockReposter creates a new mock instance.
func NewMockReposter(ctrl *gomock.Controller) *MockReposter {
	mock := &MockReposter{ctrl: ctrl}
	mock.recorder = &MockReposterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReposter) EXPECT() *MockReposterMockRecorder {
	return m.recorder
}

// Repost mocks base method.
func (m *MockReposter) Repost(ctx context.Context, user, parentID, quote string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Repost", ctx, user, parentID, quote)
	ret0, _ := ret[0].(error)
	return ret0
}

// Repost indicates an expected call of Repost.
func (mr *MockReposterMockRecorder) Repost(ctx, user, parentID, quote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Repost", reflect.TypeOf((*MockReposter)(nil).Repost), ctx, user, parentID, quote)
}
