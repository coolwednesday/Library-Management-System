package users

import (
	"SimpleRESTApi/models"
	"database/sql"
	"errors"
	// connecting sql driver.
	_ "github.com/go-sql-driver/mysql"
)

// handler defines a struct that has a UserStorer interface.
type service struct {
	store UserStorer
}

// Factory Pattern
func New(usvc UserStorer) *service {
	return &service{usvc}
}

// Handler - handle the request - 1. parse request 2. send for processing 3. send response
// Service - Do business logic (here,check if user exists by calling store)
// Store - Handle database(Insert record, get record)

// Add receives json object of user details and calls add function to add the record in the database.
func (s *service) Add(userid int, name string) error {

	if userid == 0 {
		return errors.New("user details required")

	}

	if userid/1000 < 1 || userid/1000 >= 10 {
		return errors.New("enter valid id. Must be 4 digits only")
	}

	err := s.store.CheckUser(userid)

	if err == nil {
		err := s.store.UpdateUser(userid)
		if err != nil {
			return err
		}
		return nil
	}
	err = s.store.Add(userid, name)
	if err != nil {
		return errors.New("duplicate id. User already exist. Try again")
	}

	return nil
}

// Remove receives user id and calls remove function to remove the record from the database.
func (s *service) Remove(userid int) (string, error) {

	if userid/1000 < 1 || userid/1000 >= 10 {
		return "400", errors.New("enter valid id. Must be 4 digits only")
	}

	result, err := s.store.Remove(userid)
	if err != nil {
		return "", errors.New("user cannot be removed. User must return the book before being removed")
	}

	val, err := result.RowsAffected()
	if err != nil {
		return "", err
	} else if val == 0 {
		return "404", errors.New("user with this id does not exist")
	}

	return "User removed successfully", err
}

// List receives user id and calls list function to list the record from the database.
func (s *service) List(userid int) (*models.User, error) {

	if userid/1000 < 1 || userid/1000 >= 10 {
		return nil, errors.New("enter valid id. Must be 4 digits only")
	}

	user, err := s.store.List(userid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil,
			errors.New("user with this id does not exist")
	}

	return user, err
}

// ListAll calls the listall function to fetch all the records from the database.
func (s *service) ListAll() ([]*models.User, error) {

	user := make([]*models.User, 0)

	user, err := s.store.ListAll()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no users available")
		}
		return nil, err
	}

	return user, nil
}
