package books

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var S *sql.DB

// BookHandler struct has a BookStorer interface.
type BookHandler struct {
	BookStore BookStorer
}

// RecordHandler struct has a RecordStorer interface
type RecordHandler struct {
	RentStore RecordStorer
}

type RentStore struct {
}

// Book Structure.
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Isbn   int    `json:"isbn"`
}

// BookStore struct
type BookStore struct {
}

// lendingRecord Book rented and User Mapping with isbn of book as key and user's id as value.
type lendingRecord struct {
	UserID int `json:"userid"`
	ISBN   int `json:"isbn"`
}

// Add receives a json object of book details and calls the add method defined on BookStore to add a record in the database.
func (bh *BookHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var book Book

	err = json.Unmarshal(b, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}

	if book.Isbn == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte(fmt.Sprintf("error: %v", errors.New("book details required")))))

		return
	}

	err = bh.BookStore.add(book.Isbn, book.Title, book.Author)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte(fmt.Sprintf("error: %v", err))))

		return
	}

	w.WriteHeader(http.StatusCreated)

	val := struct {
		Message string
		Isbn    int
	}{"book added successfully", book.Isbn}

	v, _ := json.Marshal(val)
	log.Println(w.Write(v))
}

// Remove receives the book's isbn and calls the remove method defined on BookStore to removes the record from the database.
func (bh *BookHandler) Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error

	// extracting the id from the request url
	id := mux.Vars(r)["isbn"]

	isbn, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: isbn must be an integer")))

		return
	}

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: enter valid isbn. Must be 5 digits only")))

		return
	}

	val, err := bh.BookStore.remove(isbn)
	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
}

// List receives book's isbn in request and calls the list method defined on BookStore and returns book details.
func (bh *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extracting the id from the request url
	id := mux.Vars(r)["isbn"]

	isbn, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: isbn must be an integer")))

		return
	}

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: enter valid isbn. Must be 5 digits only")))

		return
	}

	book, err := bh.BookStore.list(isbn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	val, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(val)
}

// ListAvailible returns the list of availible books that can be rented.
func (bh *BookHandler) ListAvailible(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := bh.BookStore.listavailible()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	val, err := json.Marshal(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(val)
}

// Borrow receives user id and book's isbn as request and returns a update message and error if any
// Adds the record under lendingRecords table in database.
func (rh *RecordHandler) Borrow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var l lendingRecord

	if err := json.Unmarshal(b, &l); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error unmarshalling lending record: %v", err)))

		return
	}

	if l.ISBN == 0 && l.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: empty request found. Enter the isbn and userid"))

		return
	} else if l.ISBN == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: isbn is missing. Try again"))

		return
	} else if l.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: userid is missing. Try again"))

		return
	}

	if l.ISBN/10000 < 1 || l.ISBN/10000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: enter valid isbn. Must be 5 digits only"))

		return
	}

	if l.UserID/1000 < 1 || l.UserID/1000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: enter valid id. Must be 4 digits only"))

		return
	}
	fmt.Println("Hello")

	val, err := rh.RentStore.borrow(l.UserID, l.ISBN)

	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

			return
		} else if val == "500" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
}

// Return receives isbn of book  as request and returns an update message and error if any
// Also removes record book's isbn from the lendingRecords table in database.
func (rh *RecordHandler) Return(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["isbn"]

	var err error

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: empty request found. Enter the isbn"))

		return
	}

	isbn, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: isbn must be an integer"))

		return
	}

	if isbn/10000 < 1 || isbn/10000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: enter valid isbn. Must be 5 digits only"))

		return
	}

	val, err := rh.RentStore.returnbook(isbn)
	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
		} else if val == "500" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
}

// add method defined on BookStore.
func (b *BookStore) add(isbn int, title string, author string) error {
	_, err := S.Exec(`INSERT INTO books (title,author,isbn) VALUES (?,?,?)`, title, author, isbn)
	if err != nil {
		return errors.New("duplicate isbn. Book already exists. Try again")

	}

	return nil
}

// remove method defined on BookStore.
func (b *BookStore) remove(isbn int) (string, error) {

	result, err := S.Exec(`DELETE FROM books WHERE isbn=?`, isbn)
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

// list method defined on BookStore.
func (b *BookStore) list(isbn int) (*Book, error) {

	var title, author string

	err := S.QueryRow(`Select title,author from books where isbn=?`, isbn).Scan(&title, &author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book with this isbn does not exist")
		}

		return nil, err
	}

	book := &Book{
		title,
		author,
		isbn,
	}
	return book, nil
}

// listavailible method defined on BookStore.
func (bs *BookStore) listavailible() ([]*Book, error) {
	books := make([]*Book, 0)

	rows, err := S.Query(`Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on 
    s.isbn = r.bookid where r.bookid is null;`)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return []*Book{}, errors.New("no books available")
		}

		return nil, err
	}

	b := &Book{}

	for rows.Next() {
		err := rows.Scan(&b.Isbn, &b.Title, &b.Author)

		if err != nil {
			return []*Book{}, err
		}

		books = append(books, &Book{
			Isbn:   b.Isbn,
			Title:  b.Title,
			Author: b.Author,
		})
	}

	return books, nil
}

// borrow method defined on BookStore.
func (r *RentStore) borrow(id, isbn int) (string, error) {
	err := S.QueryRow(`Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid
              where s.isbn=? and r.bookid is null;`, isbn).Scan(&isbn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "404", errors.New("book with this isbn does not exist or is already borrowed")
		}

		return "500", err
	}
	_, err = S.Exec(`INSERT INTO lendingRecords(userid,bookid)VALUES((Select id from users where id=?),
                                                (Select isbn from books where isbn=?));`, id, isbn)
	if err != nil {
		return "500", fmt.Errorf("borrow book event failed. Try again: %s", err)
	}

	return "book borrowed successfully", nil
}

// returnbook method defined on BookStore.
func (r *RentStore) returnbook(isbn int) (string, error) {
	result, err := S.Exec(`DELETE FROM lendingRecords where bookid=?`, isbn)
	if err != nil {
		return "500", fmt.Errorf("return book event failed. Try again: %s", err)
	}

	if val, _ := result.RowsAffected(); val == 0 {
		err := S.QueryRow(`Select * from books where isbn=?`, isbn).Scan()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "404", errors.New("book with this isbn does not exist")
			}
		}
		return "", errors.New("book with this isbn was not borrowed")
	}

	return "book returned successfully", nil
}
