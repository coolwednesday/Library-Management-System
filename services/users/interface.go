package users

import (
	"SimpleRESTApi/models"
	"database/sql"
)

type UserStorer interface {
	Add(int, string) error
	List(int) (*models.User, error)
	ListAll() ([]*models.User, error)
	CheckUser(int) error
	Remove(int) (sql.Result, error)
	UpdateUser(int) error
}
