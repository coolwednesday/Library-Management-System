package books

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	mock "github.com/libraryManagementSystem/handler"
	"github.com/libraryManagementSystem/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"net/http"
	//"net/http/httptest"
	"strconv"
	//"strings"
	"testing"
)

// TestAddBooks function tests for all possible requests to Add.
func TestAddBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		body         string
		arg1         int
		arg2         string
		arg3         string
		expectedCode int
		expectedBody res
	}{
		{"Adding a Book", "/book", `{"title":"Book6","author":"Author5","isbn":34567}`, 34567, "Book6", "Author5",
			201,
			res{
				struct {
					Message string
					Isbn    int
				}{"book added successfully", 34567},
				nil}},

		{"Empty request", "/book", `{}`, 0, "", "", 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"book details required"}}}},

		{"Adding Book with Duplicate ISBN", "/book",
			`{"title":"Book1","author":"Author5","isbn":12905}`, 12905, "Book1", "Author5", 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: duplicate isbn. Book already exists. Try again"}}}},
	}

	//Running for all testcases
	for _, test := range tests {

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, test.target, bytes.NewBuffer([]byte(test.body)))
		req.Header.Set("Content-Type", "application/json")
		request := gofrHttp.NewRequest(req)
		//response := http.NewResponder(httptest.NewRecorder(), test.input)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		mockStore.EXPECT().Add(ctx, test.arg1, test.arg2, test.arg3).Return(test.expectedBody.err)

		uh := handler{mockStore}

		val, err := uh.Add(ctx)

		res := res{val, err}
		assert.Equal(t, test.expectedBody, res)

	}
}

// TestRemoveBook function tests for all possible requests to Remove.
func TestRemoveBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		arg1         int
		expectedBody res
		expectedCode int
	}{
		{"Wrong format of ISBN, 6 Digits", "/book/345679", 345679,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid isbn. Must be 5 digits only"}}}, 400},

		{"Empty ISBN", "/book/", -1,
			res{
				nil,
				gofrHttp.ErrorMissingParam{Params: []string{"isbn"}}}, 400},

		{"ISBN does not exist", "/book/32456", 32456,
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Book", "error: book with this isbn does not exist"}}, 404},

		{"Correct ISBN", "/book/34567", 34567,
			res{
				nil,
				nil}, 204},

		{"Sending String instead of Integer", "/book/`345679`", 0,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: isbn must be an integer"}}}, 400},
	}

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, test.target, bytes.NewBuffer([]byte(``)))

		if test.arg1 != -1 {
			req = mux.SetURLVars(req, map[string]string{"isbn": strconv.Itoa(test.arg1)})
		} else if test.arg1 == -1 {
			req = mux.SetURLVars(req, map[string]string{"isbn": ""})
		} else {
			req = mux.SetURLVars(req, map[string]string{"isbn": "`345679`"})
		}

		request := gofrHttp.NewRequest(req)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		if test.arg1 != -1 {
			mockStore.EXPECT().Remove(ctx, test.arg1).Return(test.expectedBody.err)
		}

		uh := handler{mockStore}

		val, err := uh.Remove(ctx)
		res := res{val, err}
		assert.Equal(t, test.expectedBody, res, test.name)

	}
}

// TestListBook function tests for all possible requests to List.
func TestListBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		arg1         int
		expectedBody res
		expectedCode int
	}{
		{"Wrong format of ISBN, 6 Digits", "/book/345679", 345679,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid isbn. Must be 5 digits only"}}}, 400},

		{"Empty ISBN", "/book/", -1,
			res{
				nil,
				gofrHttp.ErrorMissingParam{Params: []string{"isbn"}}}, 400},

		{"ISBN does not exist", "/book/32456", 32456,
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Book", "error: book with this isbn does not exist"}}, 404},

		{"Correct ISBN", "/book/12785", 12785,
			res{
				&models.Book{"Book3", "Autho3", 12785},
				nil}, 200},

		{"Sending String instead of Integer", "/book/`345679`", 0,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: isbn must be an integer"}}}, 400},
	}

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, test.target, bytes.NewBuffer([]byte(``)))

		if test.arg1 != -1 {
			req = mux.SetURLVars(req, map[string]string{"isbn": strconv.Itoa(test.arg1)})
		} else if test.arg1 == -1 {
			req = mux.SetURLVars(req, map[string]string{"isbn": ""})
		} else {
			req = mux.SetURLVars(req, map[string]string{"isbn": "`345679`"})
		}

		request := gofrHttp.NewRequest(req)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		if test.arg1 != -1 {
			mockStore.EXPECT().List(ctx, test.arg1).Return(test.expectedBody.value, test.expectedBody.err)
		}

		uh := handler{mockStore}

		val, err := uh.List(ctx)
		res := res{val, err}
		assert.Equal(t, test.expectedBody, res, test.name)

	}
}

// TestBorrowBook function tests for all possible requests to Borrow.
func TestBorrowBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		body         string
		arg1         int
		arg2         int
		expectedCode int
		expectedBody res
	}{
		{"Empty Request", "/book/rent", "{}", 0, 0, 400,
			res{
				nil,
				gofrHttp.ErrorMissingParam{Params: []string{"empty request found. Enter the isbn and userid"}},
			}},

		{"Invalid Book", "/book/rent", `{"userid":1340,"isbn":12345}`, 1340, 12345, 404,
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Book", "book with this isbn does not exist or is already borrowed"}}},

		{"Wrong Format of UserID", "/book/rent", `{"userid":123,"isbn":12345}`, 123, 12345, 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid id. Must be 4 digits only"}},
			}},

		{"Wrong Format of ISBN", "/book/rent", `{"userid":1234,"isbn":123456}`, 1234, 123456, 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid isbn. Must be 5 digits only"}},
			}},

		{"Missing UserID", "/book/rent", `{"isbn": 34567}`, 0, 34567, 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: userid is missing. Try again."}}},
		},

		{"Missing Book ISBN", "/book/rent", `{"userid": 1234}`, 1234, 0, 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: isbn is missing. Try again."}},
			}},

		{"Available Book", "/book/rent", `{"userid":1567,"isbn":19905}`, 1567, 19905, 200,
			res{
				struct {
					Message string
					Isbn    int
				}{"book borrowed successfully", 19905},
				nil},
		},

		{"Book Already Borrowed", "/book/rent", `{"userid":1567,"isbn":19905}`, 1567, 19905, 404,
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Book", "book with this isbn does not exist or is already borrowed"},
			},
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, test.target, bytes.NewBuffer([]byte(test.body)))
		req.Header.Set("Content-Type", "application/json")
		request := gofrHttp.NewRequest(req)
		//response := http.NewResponder(httptest.NewRecorder(), test.input)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		mockStore.EXPECT().Borrow(ctx, test.arg1, test.arg2).Return(test.expectedBody.err)

		uh := handler{mockStore}

		val, err := uh.Borrow(ctx)

		res := res{val, err}
		assert.Equal(t, test.expectedBody, res)

	}
}

// TestReturnBook function tests for all possible requests to Return.
func TestReturnBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		arg1         int
		expectedBody res
		expectedCode int
	}{
		{"Empty Request", "/book/return/", -1,
			res{
				nil,
				gofrHttp.ErrorMissingParam{Params: []string{"error: empty request found. Enter the isbn"}}}, 400},

		{"Invalid Book", "/book/return/12345", 12345,
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Borrow Book Event", "error: book with this isbn does not exist"}}, 404},

		{"Wrong Format of ISBN", "/book/return/123456", 123456,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid isbn. Must be 5 digits only"}}}, 404},

		{"Book Not Borrowed", "/book/return/12905", 12905,
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Borrow Book Event", "book with this isbn was not borrowed"}}, 404},

		{"Return Book", "/book/return/19905", 19905,
			res{
				nil,
				nil}, 200},
	}

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, test.target, bytes.NewBuffer([]byte(``)))

		if test.arg1 != -1 {
			req = mux.SetURLVars(req, map[string]string{"isbn": strconv.Itoa(test.arg1)})
		} else {
			req = mux.SetURLVars(req, map[string]string{"isbn": ""})
		}

		request := gofrHttp.NewRequest(req)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		if test.arg1 != -1 {
			mockStore.EXPECT().Returnbook(ctx, test.arg1).Return(test.expectedBody.err)
		}

		uh := handler{mockStore}

		val, err := uh.Return(ctx)
		res := res{val, err}
		assert.Equal(t, test.expectedBody, res, test.name)

	}

}

// TestListBook function tests for all possible requests to TestListBookHandler.
func TestListAvailibleBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		expectedBody res
		expectedCode int
	}{
		{"No Users Availible", "/user",
			res{
				nil,
				gofrHttp.ErrorEntityNotFound{"Books", "no books availible"}}, 404},

		{
			"Users Availible", "/user",
			res{
				[]*models.Book{{"Book3", "Author3", 12785}},
				nil}, 200},
	}

	for _, test := range tests {

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, test.target, bytes.NewBuffer([]byte(``)))

		request := gofrHttp.NewRequest(req)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		mockStore.EXPECT().ListAvailible(ctx).Return(test.expectedBody.value, test.expectedBody.err)

		uh := handler{mockStore}

		val, err := uh.ListAvailible(ctx)
		res := res{val, err}
		assert.Equal(t, test.expectedBody, res, test.name)
	}
}
