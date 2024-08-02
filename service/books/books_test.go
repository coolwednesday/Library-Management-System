package books

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/libraryManagementSystem/models"
	mock "github.com/libraryManagementSystem/service"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	//"net/http/httptest"

	//"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestAddBooks function tests for all possible requests to Add.
func TestAddBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookStorer(ctrl)

	tests := []struct {
		name         string
		callsStore   int
		input        models.Book
		expectedBody error
	}{
		{"Adding a Book", 1, models.Book{"Book6", "Author5", 34567},
			nil},

		{"Empty request", -1, models.Book{},
			gofrHttp.ErrorMissingParam{[]string{"book details required"}}},

		{"Adding Book with Duplicate ISBN", 2,
			models.Book{"Book1", "Author5", 12905},
			gofrHttp.ErrorEntityAlreadyExist{}},
	}

	//Running for all testcases
	for _, test := range tests {
		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}

		bh := service{mockStore}
		if test.callsStore == 1 {
			mockStore.EXPECT().CheckBook(ctx, 34567).Return(sql.ErrNoRows)
			mockStore.EXPECT().Add(ctx, 34567, "Book6", "Author5").Return(nil)
		} else if test.callsStore == 2 {
			mockStore.EXPECT().CheckBook(ctx, 12905).Return(sql.ErrNoRows)
			mockStore.EXPECT().Add(ctx, 12905, "Book1", "Author5").Return(sql.ErrNoRows)
		}

		err := bh.Add(ctx, test.input.Isbn, test.input.Title, test.input.Author)

		assert.Equal(t, test.expectedBody, err)
	}
}

// TestRemoveBook function tests for all possible requests to Remove.
func TestRemoveBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookStorer(ctrl)

	tests := []struct {
		name         string
		callsStore   int
		input        int
		expectedBody error
	}{
		{"Wrong format of ISBN, 6 Digits", -1, 345679, gofrHttp.ErrorInvalidParam{[]string{"enter valid isbn. Must be 5 digits only"}}},
		{"ISBN does not exist", 1, 32456, gofrHttp.ErrorEntityNotFound{"Book", "books with this isbn does not exist."}},
		{"Correct ISBN", 2, 34567, nil},
	}

	//Testing all the tests for Remove.
	for _, test := range tests {
		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}
		if test.callsStore == 1 {
			mockStore.EXPECT().Remove(ctx, 32456).Return(gofrHttp.ErrorEntityNotFound{"Book", "books with this isbn does not exist."})
		} else if test.callsStore == 2 {
			mockStore.EXPECT().Remove(ctx, 34567).Return(nil)

		}

		s := service{mockStore}
		err := s.Remove(ctx, test.input)
		assert.Equal(t, test.expectedBody, err)
	}
}

// TestListBook function tests for all possible requests to List.
func TestListBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookStorer(ctrl)

	type args struct {
		books *models.Book
		err   error
	}

	tests := []struct {
		name         string
		callsStore   int
		input        int
		expectedBody args
	}{
		{"Wrong format of ISBN, 6 Digits", -1, 345679, args{nil, gofrHttp.ErrorInvalidParam{[]string{"error: enter valid isbn. Must be 5 digits only"}}}},
		{"ISBN does not exist", 1, 32456, args{nil, gofrHttp.ErrorEntityNotFound{"Book", "book with this isbn does not exist"}}},
		{"Correct ISBN", 2, 12785, args{&models.Book{"Book3", "Author3", 12785}, nil}},
	}

	for _, test := range tests {

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}

		s := service{mockStore}

		if test.callsStore == 1 {
			mockStore.EXPECT().List(ctx, 32456).Return(nil, gofrHttp.ErrorEntityNotFound{"Book", "book with this isbn does not exist"})
		} else if test.callsStore == 2 {
			mockStore.EXPECT().List(ctx, 12785).Return(&models.Book{
				"Book3",
				"Author3",
				12785,
			}, nil)

		}

		user, err := s.List(ctx, test.input)

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

	mockStore := mock.NewMockBookStorer(ctrl)

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
		callsStore   int
		body         input
		expectedBody error
	}{
		{"Empty Request", -1, input{0, 0},
			gofrHttp.ErrorMissingParam{Params: []string{"empty request found. Enter the isbn and userid"}}},
		{"Invalid Book", 1, input{1340, 12345},
			gofrHttp.ErrorEntityNotFound{"Book", "book with this isbn does not exist or is already borrowed"}},
		{"Wrong Format of UserID", -1, input{123, 12345},
			gofrHttp.ErrorInvalidParam{[]string{"enter valid id. Must be 4 digits only"}}},
		{"Wrong Format of ISBN", -1, input{1234, 123456},
			gofrHttp.ErrorInvalidParam{[]string{"enter valid isbn. Must be 5 digits only"}}},
		{"Missing UserID", -1, input{0, 34567},
			gofrHttp.ErrorMissingParam{[]string{"userid is missing. Try again"}}},
		{"Missing Book ISBN", -1, input{1234, 0},
			gofrHttp.ErrorMissingParam{[]string{"isbn is missing. Try again"}}},
		{"Available Book", 0, input{1567, 19905},
			nil},
		{"Book Already Borrowed", 1, input{1567, 19905},
			gofrHttp.ErrorEntityNotFound{"Book", "book with this isbn does not exist or is already borrowed"}},
	}

	for _, test := range tests {

		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}
		s := service{mockStore}
		if test.callsStore == 1 {
			mockStore.EXPECT().CheckAvailibleBook(ctx, test.body.isbn).Return(sql.ErrNoRows)
		} else if test.callsStore == 0 {
			mockStore.EXPECT().CheckAvailibleBook(ctx, test.body.isbn).Return(nil)
			mockStore.EXPECT().Borrow(ctx, test.body.userid, test.body.isbn).Return(nil)
		}

		err := s.Borrow(ctx, test.body.userid, test.body.isbn)

		assert.Equal(t, test.expectedBody, err)
	}
}

// TestReturnBook function tests for all possible requests to Return.
func TestReturnBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockBookStorer(ctrl)

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
		{"Invalid Book", 1, 12345, gofrHttp.ErrorEntityNotFound{"Return Book Event", "book with this isbn does not exist"}},
		{"Wrong Format of ISBN", -1, 123456, gofrHttp.ErrorInvalidParam{Params: []string{"enter valid isbn. Must be 5 digits only"}}},
		{"Book Not Borrowed", 2, 12905, gofrHttp.ErrorEntityNotFound{"Return Book Event", "book with this isbn was not borrowed"}},
		{"Return Book", 0, 19905, nil},
	}

	for _, test := range tests {

		s := service{mockStore}
		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}

		if test.callsStore == 1 {
			mockStore.EXPECT().Returnbook(ctx, 12345).Return(sqlmock.NewResult(0, 0), nil)
			mockStore.EXPECT().CheckBook(ctx, 12345).Return(sql.ErrNoRows)
		} else if test.callsStore == 0 {
			mockStore.EXPECT().Returnbook(ctx, 19905).Return(sqlmock.NewResult(1, 1), nil)
		} else if test.callsStore == 2 {
			mockStore.EXPECT().Returnbook(ctx, 12905).Return(sqlmock.NewResult(0, 0), nil)
			mockStore.EXPECT().CheckBook(ctx, 12905).Return(nil)
		}
		err := s.Returnbook(ctx, test.input)

		assert.Equal(t, test.expectedBody, err)
	}

}

func TestListAllBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		books []*models.Book
		err   error
	}

	mockStore := mock.NewMockBookStorer(ctrl)

	tests := []struct {
		name         string
		expectedBody args
	}{
		{"Error", args{nil, gofrHttp.ErrorEntityNotFound{Name: "Books", Value: "no books availible"}}},

		{"No Error", args{[]*models.Book{
			{"Book1", "Author1", 34567}}, nil}}}

	mockStore.EXPECT().ListAvailible(&gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: nil,
	}).Return(nil, sql.ErrNoRows)
	mockStore.EXPECT().ListAvailible(&gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: nil,
	}).Return([]*models.Book{
		{"Book1", "Author1", 34567}}, nil)

	for _, test := range tests {

		s := service{mockStore}
		ctx := &gofr.Context{
			Context:   context.Background(),
			Request:   nil,
			Container: nil,
		}

		val, err := s.ListAvailible(ctx)
		res := args{
			val, err,
		}
		assert.Equal(t, test.expectedBody, res)
	}
}
