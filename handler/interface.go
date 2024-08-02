package handler

import (
	"github.com/libraryManagementSystem/models"
	"gofr.dev/pkg/gofr"
)

type BookServicer interface {
	Add(c *gofr.Context, isbn int, title, author string) error
	Remove(*gofr.Context, int) error
	List(*gofr.Context, int) (*models.Book, error)
	ListAvailible(*gofr.Context) ([]*models.Book, error)
	Borrow(*gofr.Context, int, int) error
	Returnbook(*gofr.Context, int) error
}

type UserServicer interface {
	Add(*gofr.Context, int, string) error
	Remove(*gofr.Context, int) error
	List(*gofr.Context, int) (*models.User, error)
	ListAll(*gofr.Context) ([]*models.User, error)
}
