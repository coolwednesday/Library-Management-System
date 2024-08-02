package service

import (
	"database/sql"
	"github.com/libraryManagementSystem/models"
	"gofr.dev/pkg/gofr"
)

type BookStorer interface {
	CheckBook(*gofr.Context, int) error
	UpdateBook(*gofr.Context, int) error
	Add(ctx *gofr.Context, isbn int, title, author string) error
	Remove(*gofr.Context, int) error
	List(*gofr.Context, int) (*models.Book, error)
	ListAvailible(*gofr.Context) ([]*models.Book, error)
	CheckAvailibleBook(*gofr.Context, int) error
	Borrow(*gofr.Context, int, int) error
	Returnbook(*gofr.Context, int) (sql.Result, error)
}

type UserStorer interface {
	Add(*gofr.Context, int, string) error
	List(*gofr.Context, int) (*models.User, error)
	ListAll(*gofr.Context) ([]*models.User, error)
	CheckUser(*gofr.Context, int) error
	Remove(*gofr.Context, int) (sql.Result, error)
	UpdateUser(*gofr.Context, int) error
}
