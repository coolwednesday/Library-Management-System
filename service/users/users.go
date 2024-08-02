package users

import (
	"database/sql"
	"errors"
	"github.com/libraryManagementSystem/models"
	svc "github.com/libraryManagementSystem/service"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	// connecting sql driver.
	_ "github.com/go-sql-driver/mysql"
)

// handler defines a struct that has a UserStorer interface.
type service struct {
	store svc.UserStorer
}

// Factory Pattern
func New(usvc svc.UserStorer) *service {
	return &service{usvc}
}

// Add receives json object of user details and calls add function to add the record in the database.
func (s *service) Add(c *gofr.Context,userid int, name string) error {

	if userid == 0 {
		return http.ErrorInvalidParam{Params: []string{"user details required"}}
	}

	if userid/1000 < 1 || userid/1000 >= 10 {
		return http.ErrorInvalidParam{Params: []string{"enter valid id. Must be 4 digits only"}}
	}

	err := s.store.CheckUser(c,userid)

	if err == nil {
		err := s.store.UpdateUser(c,userid)
		if err != nil {
			return err
		}
		return nil
	}
	err = s.store.Add(c,userid, name)
	if err != nil {
		return http.ErrorEntityAlreadyExist{}
	}

	return nil
}

// Remove receives user id and calls remove function to remove the record from the database.
func (s *service) Remove(c *gofr.Context,userid int) (error) {

	if userid/1000 < 1 || userid/1000 >= 10 {
		return http.ErrorInvalidParam{Params: []string{"enter valid id. Must be 4 digits only"}}
	}

	result, err := s.store.Remove(c,userid)
	if err != nil {
		return http.ErrorInvalidParam{Params: []string{"user cannot be removed. User must return the book before being removed"}}
	}

	val, err := result.RowsAffected()
	if err != nil {
		return  err
	} else if val == 0 {
		return http.ErrorEntityNotFound{"User","user with this id does not exist"}
	}

	return nil
}

// List receives user id and calls list function to list the record from the database.
func (s *service) List(c *gofr.Context,userid int) (*models.User, error) {

	if userid/1000 < 1 || userid/1000 >= 10 {
		return nil, http.ErrorInvalidParam{Params: []string{"enter valid id. Must be 4 digits only"}}
	}

	user, err := s.store.List(c,userid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil,
			http.ErrorEntityNotFound{"User", "user with this id does not exist"}
	}

	return user, nil
}

// ListAll calls the listall function to fetch all the records from the database.
func (s *service) ListAll(c *gofr.Context) ([]*models.User, error) {

	user := make([]*models.User, 0)

	user, err := s.store.ListAll(c)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.ErrorEntityNotFound{"Users","no users available"}
		}
		return nil, err
	}

	return user, nil
}
