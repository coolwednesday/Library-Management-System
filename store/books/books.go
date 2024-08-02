package books

import (
	"database/sql"
	"github.com/libraryManagementSystem/models"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	//models "github.com/libraryManagementSystem/models"
	// connecting through mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// BookHandler struct has a BookStorer interface.
type store struct {
}

func New() *store {
	return &store{}
}

func (s *store) CheckBook(c *gofr.Context, isbn int) error {
	err := c.SQL.QueryRowContext(c, `Select isbn from books where isbn=? and deleted_at is not null`, isbn).Scan(&isbn)

	return err
}

func (s *store) UpdateBook(c *gofr.Context, isbn int) error {
	_, err := c.SQL.ExecContext(c, `UPDATE books set deleted_at=null where isbn=?`, isbn)

	return err
}

// add method defined on BookStore.
func (s *store) Add(c *gofr.Context, isbn int, title string, author string) error {
	_, err := c.SQL.ExecContext(c, `INSERT INTO books (title,author,isbn) VALUES (?,?,?)`, title, author, isbn)

	if err != nil {
		return err
	}
	return nil
}

// remove method defined on BookStore.
func (s *store) Remove(c *gofr.Context, isbn int) error {
	result, err := c.SQL.ExecContext(c, `UPDATE books SET deleted_at = now() WHERE isbn = ? and deleted_at is null`, isbn)
	if err != nil {
		return err

	}
	val, err := result.RowsAffected()
	if err != nil {
		return err
	} else if val == 0 {
		return http.ErrorEntityNotFound{"Book", "books with this isbn does not exist."}

	}

	return nil
}

// list method defined on BookStore.
func (s *store) List(c *gofr.Context, isbn int) (*models.Book, error) {

	var title, author string

	err := c.SQL.QueryRowContext(c, `Select title,author from books where isbn=? and deleted_at is null`, isbn).Scan(&title, &author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.ErrorEntityNotFound{"Book", "error: book with this isbn does not exist"}
		}
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
func (s *store) ListAvailible(c *gofr.Context) ([]*models.Book, error) {
	rows, err := c.SQL.QueryContext(c, `Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on`+
		` s.isbn = r.bookid where r.bookid is null and s.deleted_at is null;`)
	if err != nil {
		return nil, err
	}
	var books []*models.Book
	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.Isbn, &b.Title, &b.Author)
		if err != nil {
			return nil, err
		}
		books = append(books, &b)
	}
	return books, nil
}

func (s *store) CheckAvailibleBook(c *gofr.Context, isbn int) error {
	err := c.SQL.QueryRowContext(c, `Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid`+
		` where s.isbn=? and r.bookid is null or r.deleted_at is not null;`, isbn).Scan(&isbn)
	return err
}

// borrow method defined on BookStore.
func (s *store) Borrow(c *gofr.Context, id, isbn int) error {
	_, err := c.SQL.ExecContext(c, `INSERT INTO lendingRecords(userid,bookid)`+
		` VALUES((Select id from users where id=? and deleted_at is null),`+
		` (Select isbn from books where isbn=? and deleted_at is null));`, id, isbn)
	return err
}

// returnbook method defined on BookStore.
func (s *store) Returnbook(c *gofr.Context, isbn int) (sql.Result, error) {
	result, err := c.SQL.ExecContext(c, `UPDATE lendingRecords SET deleted_at = now()`+
		` WHERE bookid = ? and deleted_at is null`, isbn)
	return result, err
}
