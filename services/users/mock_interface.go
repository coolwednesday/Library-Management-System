// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source=interface.go -destination=mock_interface.go -package=users
//

// Package users is a generated GoMock package.
package users

import (
	models "SimpleRESTApi/models"
	sql "database/sql"
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

// Add mocks base method.
func (m *MockUserStorer) Add(arg0 int, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockUserStorerMockRecorder) Add(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUserStorer)(nil).Add), arg0, arg1)
}

// CheckUser mocks base method.
func (m *MockUserStorer) CheckUser(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockUserStorerMockRecorder) CheckUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockUserStorer)(nil).CheckUser), arg0)
}

// List mocks base method.
func (m *MockUserStorer) List(arg0 int) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockUserStorerMockRecorder) List(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUserStorer)(nil).List), arg0)
}

// ListAll mocks base method.
func (m *MockUserStorer) ListAll() ([]*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll")
	ret0, _ := ret[0].([]*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAll indicates an expected call of ListAll.
func (mr *MockUserStorerMockRecorder) ListAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockUserStorer)(nil).ListAll))
}

// Remove mocks base method.
func (m *MockUserStorer) Remove(arg0 int) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockUserStorerMockRecorder) Remove(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockUserStorer)(nil).Remove), arg0)
}

// UpdateUser mocks base method.
func (m *MockUserStorer) UpdateUser(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserStorerMockRecorder) UpdateUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserStorer)(nil).UpdateUser), arg0)
}