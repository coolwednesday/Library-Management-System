package users

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
	"strconv"
	"testing"
)

// TestAddUsers function tests for all possible requests to Add.
func TestAddUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserServicer(ctrl)

	type res struct {
		value interface{}
		err   error
	}

	tests := []struct {
		name         string
		target       string
		input        string
		arg1         int
		arg2         string
		expectedCode int
		expectedBody res
	}{
		{"Wrong format of ID, 5 Digits", "http://localhost:8080/user", `{"name":"User1","id":34567}`, 34567, "User1", 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"enter valid id. Must be 4 digits only"}}}},

		{"Adding a User", "http://localhost:8080/user", `{"name":"User3","id":1234}`, 1234, "User3", 201,
			res{
				struct {
					Message string
					Isbn    int
				}{"user added successfully", 1234},
				nil}},

		{"Empty request", "http://localhost:8080/user", `{}`, 0, "", 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: user details required"}}}},

		{"Adding User with Duplicate ID", "http://localhost:8080/user", `{"name":"User1","id":9056}`, 9056, "User1", 400,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: duplicate id. User already exist. Try again"}}}},
	}

	// Running for all testcases.
	for _, test := range tests {

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, test.target, bytes.NewBuffer([]byte(test.input)))
		req.Header.Set("Content-Type", "application/json")
		request := gofrHttp.NewRequest(req)
		//response := http.NewResponder(httptest.NewRecorder(), test.input)

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   request,
			Container: nil,
		}

		mockStore.EXPECT().Add(ctx, test.arg1, test.arg2).Return(test.expectedBody.err)

		uh := handler{mockStore}

		val, err := uh.Add(ctx)

		res := res{val, err}
		assert.Equal(t, test.expectedBody, res)

	}
}

// TestRemoveUser function tests for all possible requests to Remove.
func TestRemoveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserServicer(ctrl)

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
		{"Wrong format of ID, 5 Digits", "/user/34567", 34567,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid id. Must be 4 digits only"}}},
			400},

		{"Empty ID", "/user/", -1,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: id must be an integer"}}},
			400},
		{"ID does not exist", "/user/3245", 3245,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: user with this id does not exist"}}},
			404},

		{"ID with no books borrowed", "/user/1234", 1234,
			res{
				nil,
				nil},
			200},

		{"ID that has Borrowed a book", "/user/8902", 8902,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: user cannot be removed. User must return the book before being removed"}}},
			400},
		// {"Wrong Method", "GET", "/user?id=9056", `9056`, []byte("HTTP method \"GET\" not allowed"), 405},
		{"Sending String instead of Integer", "/user/`3456`", -1, res{
			nil,
			gofrHttp.ErrorInvalidParam{Params: []string{"error: id must be an integer"}}},
			400},
	}

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, test.target, bytes.NewBuffer([]byte(``)))

		if test.arg1 != -1 {
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(test.arg1)})
		} else {
			req = mux.SetURLVars(req, map[string]string{"id": ""})
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
		assert.Equal(t, test.expectedBody, res)

		/*
			if err != nil {
				if err != test.expectedBody.err {
					t.Errorf("%v, : Error, expected: %v, got: %v", test.name, test.expectedBody.err, err)
				}
			} else {
				if val != test.expectedBody.value {
					t.Errorf("%v, : Error, expected: %v, got: %v", test.name, test.expectedBody.value, val)
				}
			}

		*/
	}
}

// TestListBook function tests for all possible requests to TestListBookHandler.
func TestListUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserServicer(ctrl)
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
		{"Wrong format of ID, 4 Digits", "/user/34567", 34567,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: enter valid id. Must be 4 digits only"}}}, 400},

		{"Empty ISBN", "/user/", -1,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: id must be an integer"}}}, 400},

		{"ID does not exist", "/user/1234", 1234,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: user with this id does not exist"}}}, 404},

		{"Correct ID", "/user/3456", 3456,
			res{
				&models.User{"User5", 3456},
				nil}, 200},

		{"Sending String instead of Integer", "/book/`3456`", -1,
			res{
				nil,
				gofrHttp.ErrorInvalidParam{Params: []string{"error: id must be an integer"}}}, 400},
	}

	for _, test := range tests {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, test.target, bytes.NewBuffer([]byte(``)))

		if test.arg1 != -1 {
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(test.arg1)})
		} else {
			req = mux.SetURLVars(req, map[string]string{"id": ""})
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

// TestListBook function tests for all possible requests to TestListBookHandler.
func TestListAllUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserServicer(ctrl)

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
				gofrHttp.ErrorInvalidParam{Params: []string{`error: no users availible`}}}, 404},

		{"Users Availible", "/user",
			res{
				[]*models.User{{"User5", 3456}},
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

		mockStore.EXPECT().ListAll(ctx).Return(test.expectedBody.value, test.expectedBody.err)

		uh := handler{mockStore}

		val, err := uh.ListAll(ctx)
		res := res{val, err}
		assert.Equal(t, test.expectedBody, res, test.name)
	}
}
