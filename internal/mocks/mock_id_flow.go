// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/LaurenceGA/go-crev/internal/command/flow (interfaces: ConfigManipulator,Github,RepoFetcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	github "github.com/LaurenceGA/go-crev/internal/github"
	id "github.com/LaurenceGA/go-crev/internal/id"
	store "github.com/LaurenceGA/go-crev/internal/store"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConfigManipulator is a mock of ConfigManipulator interface
type MockConfigManipulator struct {
	ctrl     *gomock.Controller
	recorder *MockConfigManipulatorMockRecorder
}

// MockConfigManipulatorMockRecorder is the mock recorder for MockConfigManipulator
type MockConfigManipulatorMockRecorder struct {
	mock *MockConfigManipulator
}

// NewMockConfigManipulator creates a new mock instance
func NewMockConfigManipulator(ctrl *gomock.Controller) *MockConfigManipulator {
	mock := &MockConfigManipulator{ctrl: ctrl}
	mock.recorder = &MockConfigManipulatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigManipulator) EXPECT() *MockConfigManipulatorMockRecorder {
	return m.recorder
}

// SetCurrentID mocks base method
func (m *MockConfigManipulator) SetCurrentID(arg0 *id.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCurrentID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCurrentID indicates an expected call of SetCurrentID
func (mr *MockConfigManipulatorMockRecorder) SetCurrentID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrentID", reflect.TypeOf((*MockConfigManipulator)(nil).SetCurrentID), arg0)
}

// SetCurrentStore mocks base method
func (m *MockConfigManipulator) SetCurrentStore(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCurrentStore", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCurrentStore indicates an expected call of SetCurrentStore
func (mr *MockConfigManipulatorMockRecorder) SetCurrentStore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrentStore", reflect.TypeOf((*MockConfigManipulator)(nil).SetCurrentStore), arg0)
}

// MockGithub is a mock of Github interface
type MockGithub struct {
	ctrl     *gomock.Controller
	recorder *MockGithubMockRecorder
}

// MockGithubMockRecorder is the mock recorder for MockGithub
type MockGithubMockRecorder struct {
	mock *MockGithub
}

// NewMockGithub creates a new mock instance
func NewMockGithub(ctrl *gomock.Controller) *MockGithub {
	mock := &MockGithub{ctrl: ctrl}
	mock.recorder = &MockGithubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGithub) EXPECT() *MockGithubMockRecorder {
	return m.recorder
}

// GetRepository mocks base method
func (m *MockGithub) GetRepository(arg0 context.Context, arg1, arg2 string) (*github.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepository", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepository indicates an expected call of GetRepository
func (mr *MockGithubMockRecorder) GetRepository(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepository", reflect.TypeOf((*MockGithub)(nil).GetRepository), arg0, arg1, arg2)
}

// GetUser mocks base method
func (m *MockGithub) GetUser(arg0 context.Context, arg1 string) (*github.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*github.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockGithubMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockGithub)(nil).GetUser), arg0, arg1)
}

// MockRepoFetcher is a mock of RepoFetcher interface
type MockRepoFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockRepoFetcherMockRecorder
}

// MockRepoFetcherMockRecorder is the mock recorder for MockRepoFetcher
type MockRepoFetcherMockRecorder struct {
	mock *MockRepoFetcher
}

// NewMockRepoFetcher creates a new mock instance
func NewMockRepoFetcher(ctrl *gomock.Controller) *MockRepoFetcher {
	mock := &MockRepoFetcher{ctrl: ctrl}
	mock.recorder = &MockRepoFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepoFetcher) EXPECT() *MockRepoFetcherMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockRepoFetcher) Fetch(arg0 context.Context, arg1 string) (*store.ProofStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1)
	ret0, _ := ret[0].(*store.ProofStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockRepoFetcherMockRecorder) Fetch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockRepoFetcher)(nil).Fetch), arg0, arg1)
}
