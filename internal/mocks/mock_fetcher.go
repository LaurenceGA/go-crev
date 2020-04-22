// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/LaurenceGA/go-crev/internal/store/fetcher (interfaces: GitCloner)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	git "github.com/LaurenceGA/go-crev/internal/git"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGitCloner is a mock of GitCloner interface
type MockGitCloner struct {
	ctrl     *gomock.Controller
	recorder *MockGitClonerMockRecorder
}

// MockGitClonerMockRecorder is the mock recorder for MockGitCloner
type MockGitClonerMockRecorder struct {
	mock *MockGitCloner
}

// NewMockGitCloner creates a new mock instance
func NewMockGitCloner(ctrl *gomock.Controller) *MockGitCloner {
	mock := &MockGitCloner{ctrl: ctrl}
	mock.recorder = &MockGitClonerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitCloner) EXPECT() *MockGitClonerMockRecorder {
	return m.recorder
}

// Clone mocks base method
func (m *MockGitCloner) Clone(arg0 context.Context, arg1, arg2 string) (*git.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone", arg0, arg1, arg2)
	ret0, _ := ret[0].(*git.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Clone indicates an expected call of Clone
func (mr *MockGitClonerMockRecorder) Clone(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockGitCloner)(nil).Clone), arg0, arg1, arg2)
}
