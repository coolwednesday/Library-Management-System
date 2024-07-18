package users

import (
	"SimpleRESTApi/models"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestAddUsers function tests for all possible requests to Add.
func TestAddUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserStorer(ctrl)

	tests := []struct {
		name         string
		input        models.User
		expectedBody error
	}{
		{"Wrong format of ID, 5 Digits", models.User{"User1", 34567},
			errors.New("enter valid id. Must be 4 digits only")},

		{"Wrong format of ID, 3 Digits", models.User{"User1", 345},
			errors.New("enter valid id. Must be 4 digits only")},

		{"Adding a User", models.User{"User3", 1234},
			nil},

		{"Empty request", models.User{}, errors.New("user details required")},

		{"Adding User with Duplicate ID", models.User{"User1", 9056},
			errors.New("duplicate id. User already exist. Try again")},
	}

	mockStore.EXPECT().CheckUser(1234).Return(sql.ErrNoRows)
	mockStore.EXPECT().Add(1234, "User3").Return(nil)
	mockStore.EXPECT().CheckUser(9056).Return(sql.ErrNoRows)
	mockStore.EXPECT().Add(9056, "User1").Return(errors.New("duplicate id. User already exist. Try again"))

	// Running for all testcases.
	for _, test := range tests {

		us := service{mockStore}

		err := us.Add(test.input.Id, test.input.Name)

		assert.Equal(t, test.expectedBody, err)
	}
}

// TestRemoveUser function tests for all possible requests to Remove.
func TestRemoveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserStorer(ctrl)

	type args struct {
		message string
		err     error
	}
	tests := []struct {
		name         string
		input        int
		expectedBody args
	}{
		{"Wrong format of ID, 5 Digits", 34567, args{"400", errors.New("enter valid id. Must be 4 digits only")}},
		{"Wrong format of ID, 3 Digits", 345, args{"400", errors.New("enter valid id. Must be 4 digits only")}},
		{"ID does not exist", 3245, args{"404", errors.New("user with this id does not exist")}},
		{"ID with no books borrowed", 1234, args{"User removed successfully", nil}},
		{"ID that has Borrowed a book", 8902, args{"", errors.New("user cannot be removed. " +
			"User must return the book before being removed")}},
	}

	mockStore.EXPECT().Remove(3245).Return(sqlmock.NewResult(0, 0), nil)
	mockStore.EXPECT().Remove(1234).Return(sqlmock.NewResult(1, 1), nil)
	mockStore.EXPECT().Remove(8902).Return(nil, errors.New("user cannot be removed. User must return the book before being removed"))

	// Testing all the tests for RemoveBookHandler
	for _, test := range tests {

		s := service{mockStore}

		val, err := s.Remove(test.input)
		res := args{
			val, err,
		}

		assert.Equal(t, test.expectedBody, res)
	}
}

// TestListBook function tests for all possible requests to TestListBookHandler.
func TestListUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserStorer(ctrl)

	type args struct {
		user *models.User
		err  error
	}

	tests := []struct {
		name         string
		input        int
		expectedBody args
	}{
		{"Wrong format of ID, 4 Digits", 34567, args{nil, errors.New("enter valid id. Must be 4 digits only")}},
		{"ID does not exist", 1234, args{nil, errors.New("user with this id does not exist")}},
		{"Correct ID", 3456, args{&models.User{"User5", 3456}, nil}},
	}

	mockStore.EXPECT().List(1234).Return(nil, sql.ErrNoRows)
	mockStore.EXPECT().List(3456).Return(&models.User{"User5", 3456}, nil)

	for _, test := range tests {

		s := service{mockStore}

		val, err := s.List(test.input)
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

	mockStore := NewMockUserStorer(ctrl)

	tests := []struct {
		name         string
		expectedBody args
	}{
		{"Error", args{nil, errors.New("no users available")}},
		{"No Error", args{[]*models.User{
			{"User5", 3456}}, nil}}}

	mockStore.EXPECT().ListAll().Return(nil, sql.ErrNoRows)
	mockStore.EXPECT().ListAll().Return([]*models.User{
		{"User5", 3456}}, nil)

	for _, test := range tests {

		s := service{mockStore}

		val, err := s.ListAll()
		res := args{
			val, err,
		}
		assert.Equal(t, test.expectedBody, res)
	}
}
