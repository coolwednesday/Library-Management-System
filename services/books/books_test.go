package books

import (
	"SimpleRESTApi/models"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestAddBooks function tests for all possible requests to Add.
func TestAddBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookStorer(ctrl)

	tests := []struct {
		name         string
		input        models.Book
		expectedBody error
	}{
		{"Adding a Book", models.Book{"Book6", "Author5", 34567},
			nil},

		{"Empty request", models.Book{},
			errors.New("book details required")},

		{"Adding Book with Duplicate ISBN",
			models.Book{"Book1", "Author5", 12905},
			errors.New("duplicate isbn. Book already exists. Try again")},
	}

	mockStore.EXPECT().CheckBook(34567).Return(sql.ErrNoRows)
	mockStore.EXPECT().Add(34567, "Book6", "Author5").Return(nil)
	mockStore.EXPECT().CheckBook(12905).Return(sql.ErrNoRows)
	mockStore.EXPECT().Add(12905, "Book1", "Author5").Return(sql.ErrNoRows)

	//Running for all testcases
	for _, test := range tests {

		bh := service{mockStore}

		err := bh.Add(test.input.Isbn, test.input.Title, test.input.Author)

		assert.Equal(t, test.expectedBody, err)
	}
}

// TestRemoveBook function tests for all possible requests to Remove.
func TestRemoveBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookStorer(ctrl)

	type args struct {
		message string
		err     error
	}

	tests := []struct {
		name         string
		input        int
		expectedBody args
	}{
		{"Wrong format of ISBN, 6 Digits", 345679, args{"", errors.New("enter valid isbn. Must be 5 digits only")}},
		{"ISBN does not exist", 32456, args{"404", errors.New("book with this isbn does not exist")}},
		{"Correct ISBN", 34567, args{"book removed successfully", nil}},
	}

	mockStore.EXPECT().Remove(32456).Return(sqlmock.NewResult(0, 0), nil)
	mockStore.EXPECT().Remove(34567).Return(sqlmock.NewResult(1, 1), nil)

	//Testing all the tests for Remove.
	for _, test := range tests {

		s := service{mockStore}
		val, err := s.Remove(test.input)
		res := args{
			val, err,
		}
		assert.Equal(t, test.expectedBody, res)
	}
}

// TestListBook function tests for all possible requests to List.
func TestListBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookStorer(ctrl)

	type args struct {
		books *models.Book
		err   error
	}

	tests := []struct {
		name         string
		input        int
		expectedBody args
	}{
		{"Wrong format of ISBN, 6 Digits", 345679, args{nil, errors.New("enter valid isbn. Must be 5 digits only")}},
		{"ISBN does not exist", 32456, args{nil, errors.New("book with this isbn does not exist")}},
		{"Correct ISBN", 12785, args{&models.Book{"Book3", "Author3", 12785}, nil}},
	}

	mockStore.EXPECT().List(32456).Return(nil, errors.New("book with this isbn does not exist"))
	mockStore.EXPECT().List(12785).Return(&models.Book{
		"Book3",
		"Author3",
		12785,
	}, nil)

	for _, test := range tests {

		s := service{mockStore}

		user, err := s.List(test.input)

		res := args{
			user, err,
		}
		assert.Equal(t, test.expectedBody, res)
	}
}

// TestBorrowBook function tests for all possible requests to Borrow.
func TestBorrowBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookStorer(ctrl)

	type args struct {
		message string
		err     error
	}
	type input struct {
		userid int
		isbn   int
	}

	tests := []struct {
		name         string
		input        input
		expectedBody args
	}{
		{"Empty Request", input{0, 0},
			args{"400", errors.New("empty request found. Enter the isbn and userid")}},
		{"Invalid Book", input{1340, 12345},
			args{"404", errors.New("book with this isbn does not exist or is already borrowed")}},
		{"Wrong Format of UserID", input{123, 12345},
			args{"400", errors.New("enter valid id. Must be 4 digits only")}},
		{"Wrong Format of ISBN", input{1234, 123456},
			args{"400", errors.New("enter valid isbn. Must be 5 digits only")}},
		{"Missing UserID", input{0, 34567},
			args{"400", errors.New("userid is missing. Try again")}},
		{"Missing Book ISBN", input{1234, 0},
			args{"400", errors.New("isbn is missing. Try again")}},
		{"Available Book", input{1567, 19905},
			args{"book borrowed successfully", nil}},
		{"Book Already Borrowed", input{1567, 19905},
			args{"404", errors.New("book with this isbn does not exist or is already borrowed")}},
	}

	mockStore.EXPECT().CheckAvailibleBook(12345).Return(sql.ErrNoRows)
	mockStore.EXPECT().CheckAvailibleBook(19905).Return(nil)
	mockStore.EXPECT().Borrow(1567, 19905).Return(nil)
	mockStore.EXPECT().CheckAvailibleBook(19905).Return(sql.ErrNoRows)

	for _, test := range tests {

		s := service{mockStore}

		val, err := s.Borrow(test.input.userid, test.input.isbn)

		res := args{
			val, err,
		}

		assert.Equal(t, test.expectedBody, res)
	}
}

// TestReturnBook function tests for all possible requests to Return.
func TestReturnBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookStorer(ctrl)

	type args struct {
		message string
		err     error
	}

	tests := []struct {
		name         string
		input        int
		expectedBody args
	}{
		{"Invalid Book", 12345, args{"404", errors.New("book with this isbn does not exist")}},
		{"Wrong Format of ISBN", 123456, args{"", errors.New("enter valid isbn. Must be 5 digits only")}},
		{"Book Not Borrowed", 12905, args{"", errors.New("book with this isbn was not borrowed")}},
		{"Return Book", 19905, args{"book returned successfully", nil}},
	}

	mockStore.EXPECT().Returnbook(12345).Return(sqlmock.NewResult(0, 0), nil)
	mockStore.EXPECT().CheckBook(12345).Return(sql.ErrNoRows)
	mockStore.EXPECT().Returnbook(12905).Return(sqlmock.NewResult(0, 0), nil)
	mockStore.EXPECT().CheckBook(12905).Return(nil)
	mockStore.EXPECT().Returnbook(19905).Return(sqlmock.NewResult(1, 1), nil)

	for _, test := range tests {

		s := service{mockStore}

		val, err := s.Returnbook(test.input)

		res := args{
			val, err,
		}

		assert.Equal(t, test.expectedBody, res)
	}

}

/*
func TestListAllBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockBookStorer(ctrl)


}
*/
