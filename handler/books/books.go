package books

import (
	bh "github.com/libraryManagementSystem/handler"
	"github.com/libraryManagementSystem/models"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
	"strconv"

	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// BookHandler struct has a BookStorer interface.
type handler struct {
	service bh.BookServicer
}

func New(bs bh.BookServicer) *handler {
	return &handler{bs}
}

// Add receives a json object of book details and calls the add method defined on BookStore to add a record in the database.
func (bh *handler) Add(ctx *gofr.Context) (interface{}, error) {
	var book models.Book

	if err := ctx.Bind(&book); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, http.ErrorInvalidParam{Params: []string{"body"}}
	}

	err := bh.service.Add(ctx, book.Isbn, book.Title, book.Author)

	if err != nil {
		return nil, err
	}

	val := struct {
		Message string
		Isbn    int
	}{"book added successfully", book.Isbn}

	return val, nil
}

// Remove receives the book's isbn and calls the remove method defined on BookStore to removes the record from the database.
func (bh *handler) Remove(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("isbn")

	if id == "" {
		return nil, http.ErrorMissingParam{Params: []string{"isbn"}}
	}

	isbn, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"error: isbn must be an integer"}}
	}

	err = bh.service.Remove(ctx, isbn)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// List receives book's isbn in request and calls the list method defined on BookStore and returns book details.
func (bh *handler) List(ctx *gofr.Context) (interface{}, error) {

	isbn := ctx.PathParam("isbn")

	if isbn == "" {
		return nil, http.ErrorMissingParam{Params: []string{"isbn"}}
	}

	id, err := strconv.Atoi(isbn)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"isbn"}}
	}

	book, err := bh.service.List(ctx, id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// ListAvailible returns the list of availible books that can be rented.
func (bh *handler) ListAvailible(ctx *gofr.Context) (interface{}, error) {

	books, err := bh.service.ListAvailible(ctx)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// Borrow receives user id and book's isbn as request and returns a update message and error if any
// Adds the record under lendingRecords table in database.
func (rh *handler) Borrow(ctx *gofr.Context) (interface{}, error) {

	var l models.LendingRecord

	if err := ctx.Bind(&l); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, http.ErrorInvalidParam{Params: []string{"body"}}
	}

	err := rh.service.Borrow(ctx, l.UserID, l.ISBN)

	if err != nil {
		return nil, err
	}

	val := struct {
		Message string
		Isbn    int
	}{"book borrowed successfully", l.ISBN}

	return val, nil

}

// Return receives isbn of book  as request and returns an update message and error if any
// Also removes record book's isbn from the lendingRecords table in database.

func (rh *handler) Return(ctx *gofr.Context) (interface{}, error) {

	id := ctx.PathParam("isbn")

	if id == "" {
		return nil, http.ErrorMissingParam{Params: []string{"error: empty request found. Enter the isbn"}}
	}

	isbn, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"isbn"}}
	}

	err = rh.service.Returnbook(ctx, isbn)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
