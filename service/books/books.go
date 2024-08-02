package books

import (
	"database/sql"
	"errors"
	"github.com/libraryManagementSystem/models"
	svc "github.com/libraryManagementSystem/service"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
)

// BookHandler struct has a BookStorer interface.
type service struct {
	store svc.BookStorer
}

func New(bs svc.BookStorer) *service {
	return &service{bs}
}

// Add receives a json object of book details and calls the add method defined on BookStore to add a record in the database.
func (s *service) Add(c *gofr.Context, isbn int, title, author string) error {
	if isbn == 0 {
		return http.ErrorMissingParam{Params: []string{"book details required"}}
	}

	err := s.store.CheckBook(c, isbn)
	if err == nil {
		err = s.store.UpdateBook(c, isbn)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.store.Add(c, isbn, title, author)
	if err != nil {
		return http.ErrorEntityAlreadyExist{}
	}

	return nil
}

// Remove receives the book's isbn and calls the remove method defined on BookStore to removes the record from the database.
func (s *service) Remove(c *gofr.Context, isbn int) error {

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return http.ErrorInvalidParam{Params: []string{"enter valid isbn. Must be 5 digits only"}}
	}

	err := s.store.Remove(c, isbn)
	if err != nil {
		return err
	}

	return nil
}

// List receives book's isbn in request and calls the list method defined on BookStore and returns book details.
func (s *service) List(c *gofr.Context, isbn int) (*models.Book, error) {

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return nil, http.ErrorInvalidParam{Params: []string{"error: enter valid isbn. Must be 5 digits only"}}
	}

	book, err := s.store.List(c, isbn)
	return book, err
}

// ListAvailible returns the list of availible books that can be rented.
func (s *service) ListAvailible(c *gofr.Context) ([]*models.Book, error) {

	books, err := s.store.ListAvailible(c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.ErrorEntityNotFound{"Books", "no books availible"}
		}
		return nil, err
	}

	return books, nil
}

// Borrow receives user id and book's isbn as request and returns a update message and error if any
// Adds the record under lendingRecords table in database.
func (s *service) Borrow(c *gofr.Context, userid, isbn int) error {

	if isbn == 0 && userid == 0 {
		return http.ErrorMissingParam{Params: []string{"empty request found. Enter the isbn and userid"}}
	} else if isbn == 0 {
		return http.ErrorMissingParam{Params: []string{"isbn is missing. Try again"}}
	} else if userid == 0 {

		return http.ErrorMissingParam{Params: []string{"userid is missing. Try again"}}
	}

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return http.ErrorInvalidParam{Params: []string{"enter valid isbn. Must be 5 digits only"}}
	}

	if userid/1000 < 1 || userid/1000 >= 10 {
		return http.ErrorInvalidParam{Params: []string{"enter valid id. Must be 4 digits only"}}
	}

	err := s.store.CheckAvailibleBook(c, isbn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return http.ErrorEntityNotFound{"Book", "book with this isbn does not exist or is already borrowed"}
		}
		return err
	}

	err = s.store.Borrow(c, userid, isbn)
	if err != nil {
		return err
	}

	return nil
}

// Return receives isbn of book  as request and returns an update message and error if any
// Also removes record book's isbn from the lendingRecords table in database.
func (s *service) Returnbook(c *gofr.Context, isbn int) error {

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return http.ErrorInvalidParam{Params: []string{"enter valid isbn. Must be 5 digits only"}}
	}

	result, err := s.store.Returnbook(c, isbn)
	if err != nil {
		return err
	}

	if val, _ := result.RowsAffected(); val == 0 {
		err := s.store.CheckBook(c, isbn)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.ErrorEntityNotFound{"Return Book Event", "book with this isbn does not exist"}
			}
		}

		return http.ErrorEntityNotFound{"Return Book Event", "book with this isbn was not borrowed"}
	}

	return nil
}
