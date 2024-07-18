package books

import (
	"SimpleRESTApi/models"
	"database/sql"
	"fmt"

	"errors"
	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
)

// BookHandler struct has a BookStorer interface.
type service struct {
	store BookStorer
}

func New(bs BookStorer) *service {
	return &service{bs}
}

// Add receives a json object of book details and calls the add method defined on BookStore to add a record in the database.
func (s *service) Add(isbn int, title, author string) error {
	if isbn == 0 {
		return errors.New("book details required")
	}

	err := s.store.CheckBook(isbn)
	if err == nil {
		err = s.store.UpdateBook(isbn)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.store.Add(isbn, title, author)
	if err != nil {
		return errors.New("duplicate isbn. Book already exists. Try again")
	}

	return nil
}

// Remove receives the book's isbn and calls the remove method defined on BookStore to removes the record from the database.
func (s *service) Remove(isbn int) (string, error) {

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return "", errors.New("enter valid isbn. Must be 5 digits only")
	}

	result, err := s.store.Remove(isbn)
	if err != nil {
		return "", err
	}

	val, err := result.RowsAffected()
	if err != nil {
		return "", err
	} else if val == 0 {
		return "404", errors.New("book with this isbn does not exist")

	}
	return "book removed successfully", nil
}

// List receives book's isbn in request and calls the list method defined on BookStore and returns book details.
func (s *service) List(isbn int) (*models.Book, error) {

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return nil, errors.New("enter valid isbn. Must be 5 digits only")
	}

	book, err := s.store.List(isbn)
	return book, err
}

// ListAvailible returns the list of availible books that can be rented.
func (s *service) ListAvailible() ([]*models.Book, error) {
	books := make([]*models.Book, 0)
	rows, err := s.store.ListAvailible()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no books available")
		}
		return nil, err
	}

	b := &models.Book{}
	for rows.Next() {
		err := rows.Scan(&b.Isbn, &b.Title, &b.Author)

		if err != nil {
			return []*models.Book{}, err
		}

		books = append(books, &models.Book{
			Isbn:   b.Isbn,
			Title:  b.Title,
			Author: b.Author,
		})
	}

	return books, nil
}

// Borrow receives user id and book's isbn as request and returns a update message and error if any
// Adds the record under lendingRecords table in database.
func (s *service) Borrow(userid, isbn int) (string, error) {

	if isbn == 0 && userid == 0 {
		return "400", errors.New("empty request found. Enter the isbn and userid")
	} else if isbn == 0 {
		return "400", errors.New("isbn is missing. Try again")
	} else if userid == 0 {
		return "400", errors.New("userid is missing. Try again")
	}

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return "400", errors.New("enter valid isbn. Must be 5 digits only")
	}

	if userid/1000 < 1 || userid/1000 >= 10 {
		return "400", errors.New("enter valid id. Must be 4 digits only")
	}

	err := s.store.CheckAvailibleBook(isbn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "404", errors.New("book with this isbn does not exist or is already borrowed")
		}
		return "500", err
	}

	err = s.store.Borrow(userid, isbn)
	if err != nil {
		return "500", fmt.Errorf("borrow book event failed. Try again: %s", err)
	}

	return "book borrowed successfully", nil
}

// Return receives isbn of book  as request and returns an update message and error if any
// Also removes record book's isbn from the lendingRecords table in database.
func (s *service) Returnbook(isbn int) (string, error) {

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		return "", errors.New("enter valid isbn. Must be 5 digits only")
	}

	result, err := s.store.Returnbook(isbn)
	if err != nil {
		return "500",
			fmt.Errorf("return book event failed. Try again: %s", err)
	}

	if val, _ := result.RowsAffected(); val == 0 {
		err := s.store.CheckBook(isbn)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "404", errors.New("book with this isbn does not exist")
			}
		}

		return "", errors.New("book with this isbn was not borrowed")
	}

	return "book returned successfully", nil
}
