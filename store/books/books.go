package books

import (
	"SimpleRESTApi/models"
	"database/sql"
	"fmt"

	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// BookHandler struct has a BookStorer interface.
type store struct {
	DB *sql.DB
}

func New(db *sql.DB) *store {
	return &store{db}
}

func (s *store) CheckBook(isbn int) error {
	err := s.DB.QueryRow(`Select isbn from books where isbn=? and deleted_at is not null`, isbn).Scan(&isbn)
	return err
}

func (s *store) UpdateBook(isbn int) error {
	_, err := s.DB.Exec(`UPDATE books set deleted_at=null where isbn=?`, isbn)
	fmt.Println(err)
	return err
}

// add method defined on BookStore.
func (s *store) Add(isbn int, title string, author string) error {
	_, err := s.DB.Exec(`INSERT INTO books (title,author,isbn) VALUES (?,?,?)`, title, author, isbn)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

// remove method defined on BookStore.
func (s *store) Remove(isbn int) (sql.Result, error) {
	result, err := s.DB.Exec(`UPDATE books SET deleted_at = now() WHERE isbn = ? and deleted_at is null`, isbn)
	if err != nil {
		return nil, err

	}
	return result, nil
}

// list method defined on BookStore.
func (s *store) List(isbn int) (*models.Book, error) {

	var title, author string

	err := s.DB.QueryRow(`Select title,author from books where isbn=? and deleted_at is null`, isbn).Scan(&title, &author)
	if err != nil {
		return nil, err
	}

	book := &models.Book{
		title,
		author,
		isbn,
	}

	return book, nil
}

// listavailible method defined on BookStore.
func (s *store) ListAvailible() (*sql.Rows, error) {
	rows, err := s.DB.Query(`Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on 
    s.isbn = r.bookid where r.bookid is null and s.deleted_at is null;`)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *store) CheckAvailibleBook(isbn int) error {
	err := s.DB.QueryRow(`Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid
              where s.isbn=? and r.bookid is null or r.deleted_at is not null;`, isbn).Scan(&isbn)
	return err
}

// borrow method defined on BookStore.
func (s *store) Borrow(id, isbn int) error {

	_, err := s.DB.Exec(`INSERT INTO lendingRecords(userid,bookid)VALUES((Select id from users where id=? and deleted_at is null),
                                                (Select isbn from books where isbn=? and deleted_at is null));`, id, isbn)
	return err
}

// returnbook method defined on BookStore.
func (s *store) Returnbook(isbn int) (sql.Result, error) {
	result, err := s.DB.Exec(`UPDATE lendingRecords SET deleted_at = now() WHERE bookid = ? and deleted_at is null`, isbn)
	return result, err
}
