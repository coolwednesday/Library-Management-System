package users

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/libraryManagementSystem/models"
	mock "github.com/libraryManagementSystem/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"testing"
)

// TestAddUsers function tests for all possible requests to Add.
func TestAddUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserStorer(ctrl)

	tests := []struct {
		name         string
		callsStore   int
		input        models.User
		expectedBody error
	}{
		{"Wrong format of ID, 3 Digits", -1, models.User{"User1", 345},
			gofrHttp.ErrorInvalidParam{Params: []string{"enter valid id. Must be 4 digits only"}}},

		{"Adding a User", 1, models.User{"User3", 1234},
			nil},

		{"Empty request", -1, models.User{},
			gofrHttp.ErrorInvalidParam{Params: []string{"user details required"}}},

		{"Adding User with Duplicate ID", 1, models.User{"User1", 9056},
			gofrHttp.ErrorEntityAlreadyExist{}},
	}

	// Running for all testcases.
	for _, test := range tests {

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}
		if test.callsStore != -1 {
			mockStore.EXPECT().CheckUser(ctx, test.input.Id).Return(sql.ErrNoRows)
			mockStore.EXPECT().Add(ctx, test.input.Id, test.input.Name).Return(test.expectedBody)

		}
		us := service{mockStore}

		err := us.Add(ctx, test.input.Id, test.input.Name)

		assert.Equal(t, test.expectedBody, err)
	}
}

// TestRemoveUser function tests for all possible requests to Remove.
func TestRemoveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserStorer(ctrl)

	type args struct {
		message string
		err     error
	}

	tests := []struct {
		name         string
		callsStore   int
		input        int
		expectedBody error
	}{

		{"Wrong format of ID, 3 Digits", -1, 345, gofrHttp.ErrorInvalidParam{[]string{"enter valid id. Must be 4 digits only"}}},
		{"ID does not exist", 0, 3245, gofrHttp.ErrorEntityNotFound{"User", "user with this id does not exist"}},
		{"ID with no books borrowed", 1, 1234, nil},
		{"ID that has Borrowed a book", 2, 8902, gofrHttp.ErrorInvalidParam{[]string{"user cannot be removed. User must return the book before being removed"}}},
	}

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}
		if test.callsStore == 0 {
			mockStore.EXPECT().Remove(ctx, test.input).Return(sqlmock.NewResult(0, 0), nil)
		} else if test.callsStore == 1 {
			mockStore.EXPECT().Remove(ctx, test.input).Return(sqlmock.NewResult(1, 1), nil)
		} else if test.callsStore == 2 {
			mockStore.EXPECT().Remove(ctx, test.input).Return(nil,
				gofrHttp.ErrorInvalidParam{[]string{"user cannot be removed. User must return the book before being removed"}})
		}

		s := service{mockStore}
		err := s.Remove(ctx, test.input)

		assert.Equal(t, test.expectedBody, err)
	}
}

// TestListUser function
func TestListUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockUserStorer(ctrl)

	type args struct {
		user *models.User
		err  error
	}

	tests := []struct {
		name         string
		callsStore   int
		input        int
		expectedBody args
	}{
		{"Wrong format of ID, 4 Digits", -1, 34567,
			args{
				nil,
				gofrHttp.ErrorInvalidParam{[]string{"enter valid id. Must be 4 digits only"}}}},

		{"ID does not exist", 1, 1234,
			args{nil,
				gofrHttp.ErrorEntityNotFound{"User", "user with this id does not exist"}}},

		{"Correct ID", 2, 3456,
			args{
				&models.User{"User5", 3456}, nil}},
	}

	for _, test := range tests {
		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}
		if test.callsStore == 1 {
			mockStore.EXPECT().List(ctx, 1234).Return(nil, sql.ErrNoRows)
		} else if test.callsStore == 2 {
			mockStore.EXPECT().List(ctx, 3456).Return(&models.User{"User5", 3456}, nil)
		}

		s := service{mockStore}

		val, err := s.List(ctx, test.input)
		res := args{
			val, err,
		}

		assert.Equal(t, test.expectedBody, res)
	}
}

func TestListAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		users []*models.User
		err   error
	}

	mockStore := mock.NewMockUserStorer(ctrl)

	tests := []struct {
		name         string
		expectedBody args
	}{
		{"Error", args{nil, gofrHttp.ErrorEntityNotFound{Name: "Users", Value: "no users available"}}},

		{"No Error", args{[]*models.User{
			{"User5", 3456}}, nil}}}

	mockStore.EXPECT().ListAll(&gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: nil,
	}).Return(nil, sql.ErrNoRows)
	mockStore.EXPECT().ListAll(&gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: nil,
	}).Return([]*models.User{
		{"User5", 3456}}, nil)

	for _, test := range tests {

		s := service{mockStore}
		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}

		val, err := s.ListAll(ctx)
		res := args{
			val, err,
		}
		assert.Equal(t, test.expectedBody, res)
	}
}
