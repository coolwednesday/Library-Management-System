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

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/*
Books Management:
Each book has a title, author, and a unique ISBN.
Books can be added, removed, and listed.
*/

var S *sql.DB

// Book Structure
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Isbn   int    `json:"isbn"`
}

// lendingRecord Book rented and User Mapping with isbn of book as key and user's id as value
type lendingRecord struct {
	UserID int `json:"userid"`
	ISBN   int `json:"isbn"`
}

// Add receives a json object of book details and calls the AddBook Function to add a record in the database
func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book Book

	err = json.Unmarshal(b, &book)

	err = book.add()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte(fmt.Sprintf("Error: %v", err))))

		return
	}

	w.WriteHeader(http.StatusCreated)
	val := struct {
		Message string
		Isbn    int
	}{"Book added successfully", book.Isbn}
	v, _ := json.Marshal(val)
	log.Println(w.Write(v))

}

// RemoveBookHandle receives the book's isbn and calls the remove function to removes the record from the database
func Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err error
	//extracting the id from the request url
	id := mux.Vars(r)["isbn"]
	var book Book
	book.Isbn, err = strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ISBN must be an integer")))
		return
	}

	val, err := book.remove()
	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	}

}

// List receives book's isbn in request and calls the ListBook function and returns Book details
func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b := make([]byte, r.ContentLength)

	//extracting the id from the request url
	id := mux.Vars(r)["isbn"]

	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	book := &Book{}
	book.Isbn, err = strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ISBN must be an integer")))
		return
	}

	book, err = book.list()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	val, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(val)

}

// ListAvailibleBookHandler returns the list of availible books that can be rented
func ListAvailible(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b := Book{}
	books, err := b.listavailible()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	val, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(val)

}

// Borrow receives user id and book's isbn as request and returns a update message and error if any
// Add the record under lendingRecords table in database
func Borrow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var l lendingRecord
	err = json.Unmarshal(b, &l)
	if err := json.Unmarshal(b, &l); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error unmarshalling lending record: %v", err)))
		return
	}

	if l.ISBN == 0 && l.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Empty Request Found. Enter the BookISBN and UserID."))
		return
	} else if l.ISBN == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: BookISBN is missing. Try Again."))
		return
	} else if l.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: UserID is missing. Try Again."))
		return
	}

	val, err := l.borrow()
	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		} else if val == "500" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))

}

// Return receives isbn of book  as request and returns an update message and error if any
// Also removes record book's isbn from the lendingRecords table in database
func Return(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["isbn"]
	var err error

	var l lendingRecord
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Empty Request Found. Enter the BookISBN."))
		return
	}
	l.ISBN, err = strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ISBN must be an integer")))
		return
	}

	val, err := l.returnbook()
	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
		} else if val == "500" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	}

}

// addBooks Function
func (b *Book) add() error {

	if b.Isbn == 0 {
		return fmt.Errorf("Book Details Required")
	}
	_, err := S.Exec(`INSERT INTO books (title,author,isbn) VALUES (?,?,?)`, b.Title, b.Author, b.Isbn)
	if err != nil {
		return fmt.Errorf("Duplicate ISBN. Book Already exist. Try again!")
	}
	return nil
}

// RemoveBook Function
func (b *Book) remove() (string, error) {

	if b.Isbn/10000 < 1 || b.Isbn/10000 >= 10 {
		return "", fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}
	result, err := S.Exec(`DELETE FROM books WHERE isbn=?`, b.Isbn)
	if err != nil {
		return "", err
	} else {
		val, err := result.RowsAffected()
		if err != nil {
			return "", err
		} else if val == 0 {
			return "404", fmt.Errorf("Book with this ISBN does not exist.")
		}
	}

	return "Book removed successfully", nil

}

// list Function
func (b *Book) list() (*Book, error) {

	if b.Isbn/10000 < 1 || b.Isbn/10000 >= 10 {
		return &Book{}, fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}
	var f interface{}
	err := S.QueryRow(`Select * from books where isbn=?`, b.Isbn).Scan(&b.Isbn, &b.Title, &b.Author, &f)
	if err != nil {

		if err == sql.ErrNoRows {

			return &Book{}, errors.New("Book with this ISBN does not exist.")
		}
		return &Book{}, err
	}
	return b, nil
}

// listavailible Function
func (b *Book) listavailible() ([]*Book, error) {

	books := make([]*Book, 0)
	rows, err := S.Query(`Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where r.bookid is null;`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*Book{}, errors.New("No Books available.")
		}
		return nil, err
	}

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

// borrow Function
func (l lendingRecord) borrow() (string, error) {

	if l.ISBN/10000 < 1 || l.ISBN/10000 >= 10 {
		return "", fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}

	if l.UserID/1000 < 1 || l.UserID/1000 >= 10 {
		return "", fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	var isbn int
	err := S.QueryRow(`Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where s.isbn=? and r.bookid is null;`, l.ISBN).Scan(&isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			return "404", errors.New("Book with this ISBN does not exist or is already borrowed.")
		}
		return "500", err
	}

	_, err = S.Exec(`INSERT INTO lendingRecords(userid,bookid)VALUES((Select id from users where id=?),(Select isbn from books where isbn=?));`, l.UserID, l.ISBN)
	if err != nil {
		return "500", fmt.Errorf("Borrow Book Event Failed. Try Again: %s", err)
	}
	return "Book Borrowed Successfully", nil
}

// returnBook Function
func (l lendingRecord) returnbook() (string, error) {

	if l.ISBN/10000 < 1 || l.ISBN/10000 >= 10 {
		return "", fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}

	result, err := S.Exec(`DELETE FROM lendingRecords where bookid=?`, l.ISBN)
	if err != nil {
		return "500", fmt.Errorf("Return Book Event Failed. Try Again: %s", err)
	}

	if val, _ := result.RowsAffected(); val == 0 {
		err := S.QueryRow(`Select * from books where isbn=?`, l.ISBN).Scan()
		if err != nil {
			if err == sql.ErrNoRows {
				return "404", errors.New("Book with this ISBN does not exist.")
			}
		}
		return "", errors.New("Book with this ISBN was not borrowed.")
	}

	return "Book Returned Successfully", nil
}
