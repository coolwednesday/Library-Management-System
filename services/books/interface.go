package books

import (
	"SimpleRESTApi/models"
	"database/sql"
)

type BookStorer interface {
	Add(isbn int, title, author string) error
	Remove(int) (sql.Result, error)
	List(isbn int) (*models.Book, error)
	ListAvailible() (*sql.Rows, error)
	Borrow(int, int) error
	Returnbook(int) (sql.Result, error)
	CheckBook(int) error
	CheckAvailibleBook(int) error
	UpdateBook(int) error
}
