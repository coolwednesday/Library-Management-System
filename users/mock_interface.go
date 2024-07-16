// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -destination=mock_interface.go -package=users
//

// Package users is a generated GoMock package.
package users

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserStorer is a mock of UserStorer interface.
type MockUserStorer struct {
	ctrl     *gomock.Controller
	recorder *MockUserStorerMockRecorder
}

// MockUserStorerMockRecorder is the mock recorder for MockUserStorer.
type MockUserStorerMockRecorder struct {
	mock *MockUserStorer
}

// NewMockUserStorer creates a new mock instance.
func NewMockUserStorer(ctrl *gomock.Controller) *MockUserStorer {
	mock := &MockUserStorer{ctrl: ctrl}
	mock.recorder = &MockUserStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorer) EXPECT() *MockUserStorerMockRecorder {
	return m.recorder
}

// add mocks base method.
func (m *MockUserStorer) add(arg0 int, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// add indicates an expected call of add.
func (mr *MockUserStorerMockRecorder) add(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "add", reflect.TypeOf((*MockUserStorer)(nil).add), arg0, arg1)
}

// list mocks base method.
func (m *MockUserStorer) list(arg0 int) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "list", arg0)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// list indicates an expected call of list.
func (mr *MockUserStorerMockRecorder) list(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "list", reflect.TypeOf((*MockUserStorer)(nil).list), arg0)
}

// listall mocks base method.
func (m *MockUserStorer) listall() ([]*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "listall")
	ret0, _ := ret[0].([]*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// listall indicates an expected call of listall.
func (mr *MockUserStorerMockRecorder) listall() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "listall", reflect.TypeOf((*MockUserStorer)(nil).listall))
}

// remove mocks base method.
func (m *MockUserStorer) remove(arg0 int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "remove", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// remove indicates an expected call of remove.
func (mr *MockUserStorerMockRecorder) remove(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "remove", reflect.TypeOf((*MockUserStorer)(nil).remove), arg0)
}