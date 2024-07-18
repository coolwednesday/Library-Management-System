package users

import (
	"SimpleRESTApi/models"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAddUsers function tests for all possible requests to Add.
func TestAddUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedCode int
		expectedBody []byte
	}{
		{"Wrong format of ID, 5 Digits", "POST", "/user", `{"name":"User1","id":34567}`, 400,
			[]byte("error: enter valid id. Must be 4 digits only")},
		{"Wrong format of ID, 3 Digits", "POST", "/user", `{"name":"User1","id":345}`,
			400, []byte("error: enter valid id. Must be 4 digits only")},
		{"Adding a User", "POST", "/user", `{"name":"User3","id":1234}`, 201, []byte(`{"message":"User added successfully","id":1234}`)},
		{"Empty request", "POST", "/user", `{}`, 400, []byte("error: user details required")},
		{"Adding User with Duplicate ID", "POST", "/user", `{"name":"User1","id":9056}`, 400,
			[]byte("error: duplicate id. User already exist. Try again")},
	}

	mockStore.EXPECT().Add(34567, "User1").Return(errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().Add(345, "User1").Return(errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().Add(1234, "User3").Return(nil)
	mockStore.EXPECT().Add(0, "").Return(errors.New("user details required"))
	mockStore.EXPECT().Add(9056, "User1").Return(errors.New("duplicate id. User already exist. Try again"))

	// Running for all testcases.
	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		response := httptest.NewRecorder()

		uh := handler{mockStore}

		uh.Add(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestRemoveUser function tests for all possible requests to Remove.
func TestRemoveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ID, 5 Digits", "DELETE", "/user/34567", `34567`, []byte("error: enter valid id. Must be 4 digits only"), 400},
		{"Wrong format of ID, 3 Digits", "DELETE", "/user/345", `345`, []byte("error: enter valid id. Must be 4 digits only"), 400},
		{"Empty ID", "DELETE", "/user/id=", ``, []byte("error: id must be an integer"), 400},
		{"ID does not exist", "DELETE", "/user/3245", `3245`, []byte("error: user with this id does not exist"), 404},
		{"ID with no books borrowed", "DELETE", "/user/1234", `1234`, []byte("User removed successfully"), 200},
		{"ID that has Borrowed a book", "DELETE", "/user/8902", `8902`, []byte("error: user cannot be removed. " +
			"User must return the book before being removed"), 400},
		// {"Wrong Method", "GET", "/user?id=9056", `9056`, []byte("HTTP method \"GET\" not allowed"), 405},
		{"Sending String instead of Integer", "DELETE", "/user/`3456`", `"3456"`, []byte("error: id must be an integer"), 400},
	}

	mockStore.EXPECT().Remove(34567).Return("400", errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().Remove(345).Return("400", errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().Remove(3245).Return("404", errors.New("user with this id does not exist"))
	mockStore.EXPECT().Remove(1234).Return("User removed successfully", nil)
	mockStore.EXPECT().Remove(8902).Return("", errors.New("user cannot be removed. User must return the book before being removed"))

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		request = mux.SetURLVars(request, map[string]string{"id": test.input})

		response := httptest.NewRecorder()

		uh := handler{mockStore}

		uh.Remove(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestListBook function tests for all possible requests to TestListBookHandler.
func TestListUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ID, 4 Digits", "GET", "/user/34567", `34567`, []byte("error: enter valid id. Must be 4 digits only"), 400},
		{"Wrong format of ID, 4 Digits", "GET", "/user/345", `345`, []byte("error: enter valid id. Must be 4 digits only"), 400},
		{"Empty ISBN", "GET", "/user/", ``, []byte("error: id must be an integer"), 400},
		{"ID does not exist", "GET", "/user/1234", `1234`, []byte("error: user with this id does not exist"), 400},
		{"Correct ID", "GET", "/user/3456", `3456`, []byte(`{"name":"User5","id":3456}`), 200},
		// {"Wrong Method", "DELETE", "/user?id=3456", `3456`, []byte("HTTP method \"DELETE\" not allowed"), 405},
		{"Sending String instead of Integer", "GET", "/book/`3456`", `"3456"`, []byte("error: id must be an integer"), 400},
	}
	mockStore.EXPECT().List(34567).Return(nil, errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().List(345).Return(nil, errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().List(1234).Return(nil, errors.New("user with this id does not exist"))
	mockStore.EXPECT().List(3456).Return(&models.User{"User5", 3456}, nil)

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		request = mux.SetURLVars(request, map[string]string{"id": test.input})

		response := httptest.NewRecorder()

		uh := handler{mockStore}

		uh.List(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestListBook function tests for all possible requests to TestListBookHandler.
func TestListAllUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		expectedBody []byte
		expectedCode int
	}{
		{"No Users Availible", "GET", "/user", []byte(`error: no users availible`), 404},
		{"Users Availible", "GET", "/user", []byte(`[{"name":"User5","id":3456}]`), 200},
	}
	mockStore.EXPECT().ListAll().Return(nil, errors.New("no users availible"))
	mockStore.EXPECT().ListAll().Return([]*models.User{{"User5", 3456}}, nil)

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(""))

		response := httptest.NewRecorder()

		uh := handler{mockStore}

		uh.ListAll(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}
