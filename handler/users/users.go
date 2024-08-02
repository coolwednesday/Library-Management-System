package users

import (
	"fmt"
	h "github.com/libraryManagementSystem/handler"
	"github.com/libraryManagementSystem/models"
	"gofr.dev/pkg/gofr"
	http2 "gofr.dev/pkg/gofr/http"
	"strconv"
)

// handler defines a struct that has a UserStorer interface.
type handler struct {
	service h.UserServicer
}

// Factory Pattern
func New(us h.UserServicer) handler {
	return handler{us}
}

// Handler - handle the request - 1. parse request 2. send for processing 3. send response
// Service - Do business logic (here,check if user exists by calling store)
// Store - Handle database(Insert record, get record)

// Add receives json object of user details and calls add function to add the record in the database.
func (uh *handler) Add(ctx *gofr.Context) (interface{}, error) {

	var user models.User

	if err := ctx.Bind(&user); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, http2.ErrorInvalidParam{Params: []string{"Body"}}
	}
	fmt.Println(user)

	err := uh.service.Add(ctx, user.Id, user.Name)
	if err != nil {
		return nil, err
	}

	val := struct {
		Message string
		Isbn    int
	}{"user added successfully", user.Id}

	return val, nil
}

// Remove receives user id and calls remove function to remove the record from the database.
func (uh *handler) Remove(ctx *gofr.Context) (interface{}, error) {

	id := ctx.PathParam("id")

	userid, err := strconv.Atoi(id)
	if err != nil {
		return nil, http2.ErrorInvalidParam{Params: []string{"error: id must be an integer"}}
	}

	err = uh.service.Remove(ctx, userid)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// List receives user id and calls list function to list the record from the database.
func (uh *handler) List(ctx *gofr.Context) (interface{}, error) {

	id := ctx.PathParam("id")

	var err error
	var user *models.User

	userid, err := strconv.Atoi(id)

	if err != nil {
		return nil, http2.ErrorInvalidParam{Params: []string{"error: id must be an integer"}}
	}

	user, err = uh.service.List(ctx, userid)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ListAll calls the listall function to fetch all the records from the database.
func (uh *handler) ListAll(ctx *gofr.Context) (interface{}, error) {

	var users []*models.User

	users, err := uh.service.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
