package users

import "SimpleRESTApi/models"

type UserServicer interface {
	Add(int, string) error
	Remove(int) (string, error)
	List(int) (*models.User, error)
	ListAll() ([]*models.User, error)
}
