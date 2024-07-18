package books

import (
	"SimpleRESTApi/models"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// BookHandler struct has a BookStorer interface.
type handler struct {
	service BookServicer
}

func New(bs BookServicer) *handler {
	return &handler{bs}
}

// Add receives a json object of book details and calls the add method defined on BookStore to add a record in the database.
func (bh *handler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var book models.Book

	err = json.Unmarshal(b, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = bh.service.Add(book.Isbn, book.Title, book.Author)

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
func (bh *handler) Remove(w http.ResponseWriter, r *http.Request) {
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

	val, err := bh.service.Remove(isbn)
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
func (bh *handler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extracting the id from the request url
	id := mux.Vars(r)["isbn"]

	isbn, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: isbn must be an integer")))

		return
	}

	book, err := bh.service.List(isbn)
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
func (bh *handler) ListAvailible(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := bh.service.ListAvailible()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
func (rh *handler) Borrow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var l models.LendingRecord

	if err := json.Unmarshal(b, &l); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error unmarshalling lending record: %v", err)))

		return
	}

	val, err := rh.service.Borrow(l.UserID, l.ISBN)

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
func (rh *handler) Return(w http.ResponseWriter, r *http.Request) {
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

	val, err := rh.service.Returnbook(isbn)
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
