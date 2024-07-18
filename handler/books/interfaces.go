package books

import "SimpleRESTApi/models"

type BookServicer interface {
	Add(isbn int, title, author string) error
	Remove(isbn int) (string, error)
	List(isbn int) (*models.Book, error)
	ListAvailible() ([]*models.Book, error)
	Borrow(int, int) (string, error)
	Returnbook(int) (string, error)
}
