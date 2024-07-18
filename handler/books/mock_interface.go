// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -destination=mock_interface.go -package=books
//

// Package books is a generated GoMock package.
package books

import (
	models "SimpleRESTApi/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockBookServicer is a mock of BookServicer interface.
type MockBookServicer struct {
	ctrl     *gomock.Controller
	recorder *MockBookServicerMockRecorder
}

// MockBookServicerMockRecorder is the mock recorder for MockBookServicer.
type MockBookServicerMockRecorder struct {
	mock *MockBookServicer
}

// NewMockBookServicer creates a new mock instance.
func NewMockBookServicer(ctrl *gomock.Controller) *MockBookServicer {
	mock := &MockBookServicer{ctrl: ctrl}
	mock.recorder = &MockBookServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBookServicer) EXPECT() *MockBookServicerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockBookServicer) Add(isbn int, title, author string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", isbn, title, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockBookServicerMockRecorder) Add(isbn, title, author any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockBookServicer)(nil).Add), isbn, title, author)
}

// Borrow mocks base method.
func (m *MockBookServicer) Borrow(arg0, arg1 int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Borrow", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Borrow indicates an expected call of Borrow.
func (mr *MockBookServicerMockRecorder) Borrow(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Borrow", reflect.TypeOf((*MockBookServicer)(nil).Borrow), arg0, arg1)
}

// List mocks base method.
func (m *MockBookServicer) List(isbn int) (*models.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", isbn)
	ret0, _ := ret[0].(*models.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockBookServicerMockRecorder) List(isbn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBookServicer)(nil).List), isbn)
}

// ListAvailible mocks base method.
func (m *MockBookServicer) ListAvailible() ([]*models.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAvailible")
	ret0, _ := ret[0].([]*models.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAvailible indicates an expected call of ListAvailible.
func (mr *MockBookServicerMockRecorder) ListAvailible() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAvailible", reflect.TypeOf((*MockBookServicer)(nil).ListAvailible))
}

// Remove mocks base method.
func (m *MockBookServicer) Remove(isbn int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", isbn)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockBookServicerMockRecorder) Remove(isbn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockBookServicer)(nil).Remove), isbn)
}

// Returnbook mocks base method.
func (m *MockBookServicer) Returnbook(arg0 int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Returnbook", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Returnbook indicates an expected call of Returnbook.
func (mr *MockBookServicerMockRecorder) Returnbook(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Returnbook", reflect.TypeOf((*MockBookServicer)(nil).Returnbook), arg0)
}
