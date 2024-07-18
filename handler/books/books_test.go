package books

import (
	"SimpleRESTApi/models"
	"errors"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// TestAddBooks function tests for all possible requests to Add.
func TestAddBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedCode int
		expectedBody []byte
	}{
		{"Adding a Book", "POST", "/book", `{"title":"Book6","author":"Author5","isbn":34567}`,
			201, []byte(`{"Message":"book added successfully","Isbn":34567}`)},
		{"Empty request", "POST", "/book", `{}`, 400,
			[]byte("error: book details required")},
		{"Adding Book with Duplicate ISBN", "POST", "/book",
			`{"title":"Book1","author":"Author5","isbn":12905}`, 400,
			[]byte("error: duplicate isbn. Book already exists. Try again")},
	}

	mockStore.EXPECT().Add(34567, "Book6", "Author5").
		Return(nil)
	mockStore.EXPECT().Add(0, "", "").
		Return(errors.New("book details required"))
	mockStore.EXPECT().Add(12905, "Book1", "Author5").
		Return(errors.New("duplicate isbn. Book already exists. Try again"))

	//Running for all testcases
	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		response := httptest.NewRecorder()

		bh := handler{mockStore}

		bh.Add(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected %v, got %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestRemoveBook function tests for all possible requests to Remove.
func TestRemoveBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ISBN, 6 Digits", "DELETE", "/book/345679", `345679`, []byte("error: enter valid isbn. Must be 5 digits only"), 400},
		{"Wrong format of ISBN, 4 Digits", "DELETE", "/book/3456", `3456`, []byte("error: enter valid isbn. Must be 5 digits only"), 400},
		{"Empty ISBN", "DELETE", "/book", ``, []byte("error: isbn must be an integer"), 400},
		{"ISBN does not exist", "DELETE", "/book/32456", `32456`, []byte("error: book with this isbn does not exist"), 404},
		{"Correct ISBN", "DELETE", "/book/34567", `34567`, []byte("book removed successfully"), 200},
		{"Sending String instead of Integer", "DELETE", "/book/`345679`", `"345679"`, []byte("error: isbn must be an integer"), 400},
	}
	mockStore.EXPECT().Remove(345679).Return("400", errors.New("enter valid isbn. Must be 5 digits only"))
	mockStore.EXPECT().Remove(3456).Return("400", errors.New("enter valid isbn. Must be 5 digits only"))
	mockStore.EXPECT().Remove(32456).Return("404", errors.New("book with this isbn does not exist"))
	mockStore.EXPECT().Remove(34567).Return("book removed successfully", nil)

	//Testing all the tests for Remove.
	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))
		request = mux.SetURLVars(request, map[string]string{"isbn": test.input})
		response := httptest.NewRecorder()
		bh := handler{mockStore}
		bh.Remove(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestListBook function tests for all possible requests to List.
func TestListBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Wrong format of ISBN, 6 Digits", "GET", "/book/345679", `345679`, []byte("error: enter valid isbn. Must be 5 digits only"), 400},
		{"Wrong format of ISBN, 4 Digits", "GET", "/book/3456", `3456`, []byte("error: enter valid isbn. Must be 5 digits only"), 400},
		{"Empty ISBN", "GET", "/book/", ``, []byte("error: isbn must be an integer"), 400},
		{"ISBN does not exist", "GET", "/book/32456", `32456`, []byte("error: book with this isbn does not exist"), 400},
		{"Correct ISBN", "GET", "/book/12785", `12785`, []byte(`{"title":"Book3","author":"Author3","isbn":12785}`), 200},
		{"Sending String instead of Integer", "GET", "/book/`345679`", `"345679"`, []byte("error: isbn must be an integer"), 400},
	}
	mockStore.EXPECT().List(345679).Return(nil, errors.New("enter valid isbn. Must be 5 digits only"))
	mockStore.EXPECT().List(3456).Return(nil, errors.New("enter valid isbn. Must be 5 digits only"))
	mockStore.EXPECT().List(32456).Return(nil, errors.New("book with this isbn does not exist"))
	mockStore.EXPECT().List(12785).Return(&models.Book{
		"Book3",
		"Author3",
		12785,
	}, nil)

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		request = mux.SetURLVars(request, map[string]string{"isbn": test.input})

		response := httptest.NewRecorder()

		bh := handler{mockStore}

		bh.List(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestBorrowBook function tests for all possible requests to Borrow.
func TestBorrowBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Empty Request", "POST", "/book/rent", "{}",
			[]byte("error: empty request found. Enter the isbn and userid"), 400},
		{"Invalid Book", "POST", "/book/rent", `{"userid":1340,"isbn":12345}`,
			[]byte("error: book with this isbn does not exist or is already borrowed"), 404},
		{"Wrong Format of UserID", "POST", "/book/rent", `{"userid":123,"isbn":12345}`,
			[]byte("error: enter valid id. Must be 4 digits only"), 400},
		{"Wrong Format of ISBN", "POST", "/book/rent", `{"userid":1234,"isbn":123456}`,
			[]byte("error: enter valid isbn. Must be 5 digits only"), 400},
		{"Missing UserID", "POST", "/book/rent", `{"isbn": 34567}`,
			[]byte("error: userid is missing. Try again"), 400},
		{"Missing Book ISBN", "POST", "/book/rent", `{"userid": 1234}`,
			[]byte("error: isbn is missing. Try again"), 400},
		{"Available Book", "POST", "/book/rent", `{"userid":1567,"isbn":19905}`,
			[]byte("book borrowed successfully"), 200},
		{"Book Already Borrowed", "POST", "/book/rent", `{"userid":1567,"isbn":19905}`,
			[]byte("error: book with this isbn does not exist or is already borrowed"), 404},
	}

	mockStore.EXPECT().Borrow(0, 0).Return("400", errors.New("empty request found. Enter the isbn and userid"))
	mockStore.EXPECT().Borrow(1340, 12345).Return("404", errors.New("book with this isbn does not exist or is already borrowed"))
	mockStore.EXPECT().Borrow(123, 12345).Return("400", errors.New("enter valid id. Must be 4 digits only"))
	mockStore.EXPECT().Borrow(1234, 123456).Return("400", errors.New("enter valid isbn. Must be 5 digits only"))
	mockStore.EXPECT().Borrow(0, 34567).Return("400", errors.New("userid is missing. Try again"))
	mockStore.EXPECT().Borrow(1234, 0).Return("400", errors.New("isbn is missing. Try again"))
	mockStore.EXPECT().Borrow(1567, 19905).Return("book borrowed successfully", nil)
	mockStore.EXPECT().Borrow(1567, 19905).Return("404", errors.New("book with this isbn does not exist or is already borrowed"))

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		response := httptest.NewRecorder()

		rh := handler{mockStore}

		rh.Borrow(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}
}

// TestReturnBook function tests for all possible requests to Return.
func TestReturnBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		input        string
		expectedBody []byte
		expectedCode int
	}{
		{"Empty Request", "DELETE", "/book/return/", "", []byte("error: empty request found. Enter the isbn"), 400},
		{"Invalid Book", "DELETE", "/book/return/12345", `12345`, []byte("error: book with this isbn does not exist"), 404},
		{"Wrong Format of ISBN", "DELETE", "/book/return/123456", `123456`, []byte("error: enter valid isbn. Must be 5 digits only"), 400},
		{"Book Not Borrowed", "DELETE", "/book/return/12905", `12905`, []byte("error: book with this isbn was not borrowed"), 400},
		{"Return Book", "DELETE", "/book/return/19905", `19905`, []byte("book returned successfully"), 200},
	}

	mockStore.EXPECT().Returnbook(12345).Return("404", errors.New("book with this isbn does not exist"))
	mockStore.EXPECT().Returnbook(123456).Return("400", errors.New("enter valid isbn. Must be 5 digits only"))
	mockStore.EXPECT().Returnbook(12905).Return("", errors.New("book with this isbn was not borrowed"))
	mockStore.EXPECT().Returnbook(19905).Return("book returned successfully", nil)

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(test.input))

		request = mux.SetURLVars(request, map[string]string{"isbn": test.input})

		response := httptest.NewRecorder()

		rh := handler{mockStore}

		rh.Return(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}

}

// TestReturnBook function tests for all possible requests to Return.
func TestListAllBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookServicer(ctrl)

	tests := []struct {
		name         string
		method       string
		target       string
		expectedBody []byte
		expectedCode int
	}{
		{"No Books Availible", "GET", "/book", []byte("error: no books availible"), 404},
		{"Books Availible", "GET", "/book", []byte(`[{"title":"Book3","author":"Author3","isbn":12785}]`), 200},
	}

	mockStore.EXPECT().ListAvailible().Return(nil, errors.New("no books availible"))
	mockStore.EXPECT().ListAvailible().Return([]*models.Book{{"Book3", "Author3", 12785}}, nil)

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.target, strings.NewReader(""))

		response := httptest.NewRecorder()

		rh := handler{mockStore}

		rh.ListAvailible(response, request)

		if response.Code != test.expectedCode {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, test.expectedCode, response.Code)
		}

		if response.Body.String() != string(test.expectedBody) {
			t.Errorf("%v : Error , expected: %v, got: %v ", test.name, string(test.expectedBody), response.Body.String())
		}
	}

}