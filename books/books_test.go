package books

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAddBooks function tests for all possible requests to Add
func TestAddBooks(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedCode int
		expectedBody []byte
	}{
		{"Adding a Book", "POST", "/book", `{"title":"Book6","author":"Author5","isbn":34567}`, 201, []byte(`{"Message":"Book added successfully","Isbn":34567}`)},
		{"Empty request", "POST", "/book", `{}`, 400, []byte("Error: Book Details Required")},
		{"Adding Book with Duplicate ISBN", "POST", "/book", `{"title":"Book1","author":"Author5","isbn":12905}`, 400, []byte("Error: Duplicate ISBN. Book Already exist. Try again!")},
		//{"Wrong Method", "GET", "/book", "Divya", 405, []byte("HTTP method \"GET\" not allowed")},
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
		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestRemoveBook function tests for all possible requests to Remove
func TestRemoveBook(t *testing.T) {
	//res struct describes the output of the updated slice of Books after removing book and the error if any
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ISBN, 6 Digits", "DELETE", "/book?isbn=345679", `345679`, []byte("Error: Enter valid ISBN. Must be 5 Digits only."), 400},
		{"Wrong format of ISBN, 4 Digits", "DELETE", "/book?isbn=3456", `3456`, []byte("Error: Enter valid ISBN. Must be 5 Digits only."), 400},
		{"Empty ISBN", "DELETE", "/book?isbn=", ``, []byte("Error: ISBN must be an integer"), 400},
		{"ISBN does not exist", "DELETE", "/book?isbn=32456", `32456`, []byte("Error: Book with this ISBN does not exist."), 404},
		{"Correct ISBN", "DELETE", "/book?isbn=34567", `34567`, []byte("Book removed successfully"), 200},
		//{"Wrong Method", "GET", "/book?isbn=34567", `34567`, []byte("HTTP method \"GET\" not allowed"), 405},
		{"Sending String instead of Integer", "DELETE", "/book?isbn=`345679`", `"345679"`, []byte("Error: ISBN must be an integer"), 400},
	}

	//Testing all the tests for RemoveBook Function
	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		request = mux.SetURLVars(request, map[string]string{"isbn": test.input})
		response := httptest.NewRecorder()
		Remove(response, request)
		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}
		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestListBook function tests for all possible requests to List
func TestListBook(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ISBN, 6 Digits", "GET", "/book/345679", `345679`, []byte("Error: Enter valid ISBN. Must be 5 Digits only."), 400},
		{"Wrong format of ISBN, 4 Digits", "GET", "/book/3456", `3456`, []byte("Error: Enter valid ISBN. Must be 5 Digits only."), 400},
		{"Empty ISBN", "GET", "/book/", ``, []byte("Error: ISBN must be an integer"), 400},
		{"ISBN does not exist", "GET", "/book/32456", `32456`, []byte("Error: Book with this ISBN does not exist."), 400},
		{"Correct ISBN", "GET", "/book/12785", `12785`, []byte(`{"title":"Book3","author":"Author3","isbn":12785}`), 200},
		//{"Wrong Method", "DELETE", "/book?isbn=34567", `34567`, []byte("HTTP method \"DELETE\" not allowed"), 405},
		{"Sending String instead of Integer", "GET", "/book/`345679`", `"345679"`, []byte("Error: ISBN must be an integer"), 400},
	}

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		request = mux.SetURLVars(request, map[string]string{"isbn": test.input})
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

// TestBorrowBook function tests for all possible requests to Borrow
func TestBorrowBook(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Empty Request", "POST", "/book/rent", "{}", []byte("Error: Empty Request Found. Enter the BookISBN and UserID."), 400},
		{"Invalid Book", "POST", "/book/rent", `{"userid":1340,"isbn":12345}`, []byte("Error: Book with this ISBN does not exist or is already borrowed."), 404},
		{"Wrong Format of UserID", "POST", "/book/rent", `{"userid":123,"isbn":12345}`, []byte("Error: Enter valid ID. Must be 4 Digits only."), 400},
		{"Wrong Format of ISBN", "POST", "/book/rent", `{"userid":1234,"isbn":123456}`, []byte("Error: Enter valid ISBN. Must be 5 Digits only."), 400},
		{"Missing UserID", "POST", "/book/rent", `{"isbn": 34567}`, []byte("Error: UserID is missing. Try Again."), 400},
		{"Missing Book ISBN", "POST", "/book/rent", `{"userid": 1234}`, []byte("Error: BookISBN is missing. Try Again."), 400},
		{"Available Book", "POST", "/book/rent", `{"userid":1567,"isbn":19905}`, []byte("Book Borrowed Successfully"), 200},
		{"Book Already Borrowed", "POST", "/book/rent", `{"userid":1567,"isbn":19905}`, []byte("Error: Book with this ISBN does not exist or is already borrowed."), 404},
		//{"Wrong Method", "GET", "/book/borrow", `{"userID":9056,"bookISBN":19905}`, []byte("HTTP method \"GET\" not allowed"), 405},
	}

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		response := httptest.NewRecorder()

		Borrow(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}
		if string(response.Body.Bytes()) != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), string(response.Body.Bytes()))
		}
	}
}

// TestReturnBook function tests for all possible requests to Return
func TestReturnBook(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Empty Request", "DELETE", "/book/return/", "", []byte("Error: Empty Request Found. Enter the BookISBN."), 400},
		{"Invalid Book", "DELETE", "/book/return/12345", `12345`, []byte("Error: Book with this ISBN does not exist."), 404},
		{"Wrong Format of ISBN", "DELETE", "/book/return/123456", `123456`, []byte("Error: Enter valid ISBN. Must be 5 Digits only."), 400},
		{"Book Not Borrowed", "DELETE", "/book/return/12905", `12905`, []byte("Error: Book with this ISBN was not borrowed."), 400},
		{"Return Book", "DELETE", "/book/return/19905", `19905`, []byte("Book Returned Successfully"), 200},
		//{"Wrong Method", "POST", "/book/return", `19905`, []byte("HTTP method \"POST\" not allowed"), 405},
	}

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		request = mux.SetURLVars(request, map[string]string{"isbn": test.input})
		response := httptest.NewRecorder()

		Return(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}
		if string(response.Body.Bytes()) != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), string(response.Body.Bytes()))
		}
	}

}
