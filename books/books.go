package books

import (
	db "SimpleRESTApi/database_conn"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strconv"
)

/*
Books Management:
Each book has a title, author, and a unique ISBN.
Books can be added, removed, and listed.
*/

// Book Structure
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Isbn   int    `json:"isbn"`
}

// lendingRecord Book rented and User Mapping with isbn of book as key and user's id as value
type lendingRecord struct {
	UserID   int `json:"userid"`
	BookISBN int `json:"bookisbn"`
}

// AddBookHandler receives a json object of book details and calls the AddBook Function to add a record in the database
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book Book
	err = json.Unmarshal(b, &book)
	fmt.Println(book)
	switch r.Method {
	case http.MethodPost:
		err = book.addBook()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Book added successfully"))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// RemoveBookHandle receives the book's isbn and calls the removeBook function to removes the record from the database
func RemoveBookHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	var book Book
	book.Isbn, err = strconv.Atoi(string(b))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ISBN must be an integer")))
		return
	}

	switch r.Method {
	case http.MethodDelete:
		val, err := book.removeBook()
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// ListBookHandler receives book's isbn in request and calls the ListBook function and returns Book details
func ListBookHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	book := &Book{}
	book.Isbn, err = strconv.Atoi(string(b))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ISBN must be an integer")))
		return
	}

	switch r.Method {
	case http.MethodGet:

		book, err = book.listBook()

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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// ListAvailibleBookHandler returns the list of availible books that can be rented
func ListAvailibleBooksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello")
	switch r.Method {
	case http.MethodGet:
		b := Book{}
		books, err := b.listAvailibleBooks()
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// BorrowBookHandler receives user id and book's isbn as request and returns a update message and error if any
// Add the record under lendingRecords table in database
func BorrowBookHandler(w http.ResponseWriter, r *http.Request) {
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

	if l.BookISBN == 0 && l.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Empty Request Found. Enter the BookISBN and UserID."))
		return
	} else if l.BookISBN == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: BookISBN is missing. Try Again."))
		return
	} else if l.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: UserID is missing. Try Again."))
		return
	}

	switch r.Method {
	case http.MethodPost:
		val, err := l.borrowBook()
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}
}

// ReturnBookHandler receives isbn of book  as request and returns an update message and error if any
// Also removes record book's isbn from the lendingRecords table in database
func ReturnBookHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	var l lendingRecord
	if string(b) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Empty Request Found. Enter the BookISBN."))
		return
	}
	l.BookISBN, err = strconv.Atoi(string(b))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ISBN must be an integer")))
		return
	}

	switch r.Method {
	case http.MethodDelete:
		val, err := l.returnBook()
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// addBooks Function
func (b *Book) addBook() error {
	//Connecting to database
	s, err := db.NewConnection()
	if err != nil {
		return err
	}
	defer s.Close()

	if b.Isbn == 0 {
		return fmt.Errorf("Book Details Required")
	}
	_, err = s.Exec(`INSERT INTO books (title,author,isbn) VALUES (?,?,?)`, b.Title, b.Author, b.Isbn)
	if err != nil {
		return fmt.Errorf("Duplicate ISBN. Book Already exist. Try again!")
	}
	return nil
}

// RemoveBook Function
func (b *Book) removeBook() (string, error) {

	s, err := db.NewConnection()
	if err != nil {
		return "", err
	}
	defer s.Close()

	if b.Isbn/10000 < 1 || b.Isbn/10000 >= 10 {
		return "", fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}
	result, err := s.Exec(`DELETE FROM books WHERE isbn=?`, b.Isbn)
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

// listBook Function
func (b *Book) listBook() (*Book, error) {

	s, err := db.NewConnection()
	if err != nil {
		return &Book{}, err
	}
	defer s.Close()

	if b.Isbn/10000 < 1 || b.Isbn/10000 >= 10 {
		return &Book{}, fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}
	err = s.QueryRow(`Select * from books where isbn=?`, b.Isbn).Scan(&b.Isbn, &b.Title, &b.Author)
	if err != nil {

		if err == sql.ErrNoRows {

			return &Book{}, errors.New("Book with this ISBN does not exist.")
		}
		return &Book{}, err
	}
	return b, nil
}

// listAvailibleBooks Function
func (b *Book) listAvailibleBooks() ([]*Book, error) {
	s, err := db.NewConnection()
	if err != nil {
		return nil, err
	}
	defer s.Close()
	books := make([]*Book, 0)
	rows, err := s.Query(`Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where r.bookid is null;`)
	if err != nil {
		if err == sql.ErrNoRows {
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

// borrowBook Function
func (l lendingRecord) borrowBook() (string, error) {

	s, err := db.NewConnection()
	if err != nil {
		return "500", err
	}
	defer s.Close()

	if l.BookISBN/10000 < 1 || l.BookISBN/10000 >= 10 {
		return "", fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}

	if l.UserID/1000 < 1 || l.UserID/1000 >= 10 {
		return "", fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	var isbn int
	err = s.QueryRow(`Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where s.isbn=? and r.bookid is null;`, l.BookISBN).Scan(&isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			return "404", errors.New("Book with this ISBN does not exist or is already borrowed.")
		}
		return "500", err
	}

	_, err = s.Exec(`INSERT INTO lendingRecords(userid,bookid)VALUES((Select id from users where id=?),(Select isbn from books where isbn=?));`, l.UserID, l.BookISBN)
	if err != nil {
		return "500", fmt.Errorf("Borrow Book Event Failed. Try Again: %s", err)
	}
	return "Book Borrowed Successfully", nil
}

// returnBook Function
func (l lendingRecord) returnBook() (string, error) {

	s, err := db.NewConnection()
	if err != nil {
		return "500", err
	}
	defer s.Close()

	if l.BookISBN/10000 < 1 || l.BookISBN/10000 >= 10 {
		return "", fmt.Errorf("Enter valid ISBN. Must be 5 Digits only.")
	}
	result, err := s.Exec(`DELETE FROM lendingRecords where bookid=?`, l.BookISBN)
	if err != nil {
		return "500", fmt.Errorf("Return Book Event Failed. Try Again: %s", err)
	}
	if val, _ := result.RowsAffected(); val == 0 {
		err := s.QueryRow(`Select * from books where isbn=?`, l.BookISBN).Scan()
		if err != nil {
			if err == sql.ErrNoRows {
				return "404", errors.New("Book with this ISBN does not exist.")
			}
		}
		return "", errors.New("Book with this ISBN was not borrowed.")
	}

	return "Book Returned Successfully", nil
}
