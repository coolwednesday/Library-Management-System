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
	models "SimpleRESTApi/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserServicer is a mock of UserServicer interface.
type MockUserServicer struct {
	ctrl     *gomock.Controller
	recorder *MockUserServicerMockRecorder
}

// MockUserServicerMockRecorder is the mock recorder for MockUserServicer.
type MockUserServicerMockRecorder struct {
	mock *MockUserServicer
}

// NewMockUserServicer creates a new mock instance.
func NewMockUserServicer(ctrl *gomock.Controller) *MockUserServicer {
	mock := &MockUserServicer{ctrl: ctrl}
	mock.recorder = &MockUserServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServicer) EXPECT() *MockUserServicerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockUserServicer) Add(arg0 int, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockUserServicerMockRecorder) Add(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUserServicer)(nil).Add), arg0, arg1)
}

// List mocks base method.
func (m *MockUserServicer) List(arg0 int) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockUserServicerMockRecorder) List(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUserServicer)(nil).List), arg0)
}

// ListAll mocks base method.
func (m *MockUserServicer) ListAll() ([]*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll")
	ret0, _ := ret[0].([]*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAll indicates an expected call of ListAll.
func (mr *MockUserServicerMockRecorder) ListAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockUserServicer)(nil).ListAll))
}

// Remove mocks base method.
func (m *MockUserServicer) Remove(arg0 int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockUserServicerMockRecorder) Remove(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockUserServicer)(nil).Remove), arg0)
}
