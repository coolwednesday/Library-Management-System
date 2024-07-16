package users

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAddUsers function tests for all possible requests to Add
func TestAddUsers(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedCode int
		expectedBody []byte
	}{
		{"Wrong format of ID, 5 Digits", "POST", "/user", `{"name":"User1","id":34567}`, 400, []byte("Error: Enter valid ID. Must be 4 Digits only.")},
		{"Wrong format of ID, 3 Digits", "POST", "/user", `{"name":"User1","id":345}`, 400, []byte("Error: Enter valid ID. Must be 4 Digits only.")},
		{"Adding a User", "POST", "/user", `{"name":"User3","id":1234}`, 201, []byte(`{"message":"User added successfully","id":1234}`)},
		{"Empty request", "POST", "/user", `{}`, 400, []byte("Error: User Details Required")},
		{"Adding User with Duplicate ID", "POST", "/user", `{"title":"User1","id":9056}`, 400, []byte("Error: Duplicate ID. User Already exist. Try again!")},
		//{"Wrong Method", "GET", "/user", "Divya", 405, []byte("HTTP method \"GET\" not allowed")},
	}

	var err error
	S, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/library")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database connected", S)
	}

	//Running for all testcases
	for _, test := range tests {

		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		response := httptest.NewRecorder()
		Add(response, request)
		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, test.expectedCode, response.Code)
		}
		if string(response.Body.Bytes()) != string(test.expectedBody) {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, string(test.expectedBody), string(response.Body.Bytes()))
		}
	}
}

// TestRemoveUser function tests for all possible requests to Remove
func TestRemoveUser(t *testing.T) {

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ID, 5 Digits", "DELETE", "/user/34567", `34567`, []byte("Error: Enter valid ID. Must be 4 Digits only."), 400},
		{"Wrong format of ID, 3 Digits", "DELETE", "/user/345", `345`, []byte("Error: Enter valid ID. Must be 4 Digits only."), 400},
		{"Empty ID", "DELETE", "/user/id=", ``, []byte("Error: ID must be an integer"), 400},
		{"ID does not exist", "DELETE", "/user/3245", `3245`, []byte("Error: User with this ID does not exist."), 404},
		{"ID with no books borrowed", "DELETE", "/user/1234", `1234`, []byte("User removed successfully"), 200},
		{"ID that has Borrowed a book", "DELETE", "/user/8902", `8902`, []byte("Error: User cannot be removed. User must return the book before being removed."), 400},
		//{"Wrong Method", "GET", "/user?id=9056", `9056`, []byte("HTTP method \"GET\" not allowed"), 405},
		{"Sending String instead of Integer", "DELETE", "/user/`3456`", `"3456"`, []byte("Error: ID must be an integer"), 400},
	}

	//Testing all the tests for RemoveBookHandler
	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		request = mux.SetURLVars(request, map[string]string{"id": test.input})
		response := httptest.NewRecorder()
		Remove(response, request)
		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}
		if string(response.Body.Bytes()) != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), string(response.Body.Bytes()))
		}
	}
}

// TestListBook function tests for all possible requests to TestListBookHandler
func TestListUser(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ID, 4 Digits", "GET", "/user/34567", `34567`, []byte("Error: Enter valid ID. Must be 4 Digits only."), 400},
		{"Wrong format of ID, 4 Digits", "GET", "/user/345", `345`, []byte("Error: Enter valid ID. Must be 4 Digits only."), 400},
		{"Empty ISBN", "GET", "/user/", ``, []byte("Error: ID must be an integer"), 400},
		{"ID does not exist", "GET", "/user/1234", `1234`, []byte("Error: User with this ID does not exist."), 400},
		{"Correct ID", "GET", "/user/3456", `3456`, []byte(`{"name":"User5","id":3456}`), 200},
		//{"Wrong Method", "DELETE", "/user?id=3456", `3456`, []byte("HTTP method \"DELETE\" not allowed"), 405},
		{"Sending String instead of Integer", "GET", "/book/`3456`", `"3456"`, []byte("Error: ID must be an integer"), 400},
	}

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		request = mux.SetURLVars(request, map[string]string{"id": test.input})
		response := httptest.NewRecorder()

		List(response, request)
		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}
		if string(response.Body.Bytes()) != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), string(response.Body.Bytes()))
		}
	}

}
