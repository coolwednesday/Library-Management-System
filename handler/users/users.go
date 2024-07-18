package users

import (
	"SimpleRESTApi/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	// connecting sql driver.
	_ "github.com/go-sql-driver/mysql"
)

// handler defines a struct that has a UserStorer interface.
type handler struct {
	service UserServicer
}

// Factory Pattern
func New(us UserServicer) handler {
	return handler{us}
}

// Handler - handle the request - 1. parse request 2. send for processing 3. send response
// Service - Do business logic (here,check if user exists by calling store)
// Store - Handle database(Insert record, get record)

// Add receives json object of user details and calls add function to add the record in the database.
func (uh *handler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// b is the slice of bytes to read the request body
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	_ = json.Unmarshal(b, &user)

	err = uh.service.Add(user.Id, user.Name)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	w.WriteHeader(http.StatusCreated)
	val := struct {
		Message string `json:"message"`
		Id      int    `json:"id"`
	}{
		"User added successfully",
		user.Id,
	}
	v, _ := json.Marshal(val)
	_, _ = w.Write(v)
}

// Remove receives user id and calls remove function to remove the record from the database.
func (uh *handler) Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extracting the id from the request path
	id := mux.Vars(r)["id"]

	var err error

	userid, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: id must be an integer"))

		return
	}

	val, err := uh.service.Remove(userid)

	if err != nil {
		if val == "404" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(val))
}

// List receives user id and calls list function to list the record from the database.
func (uh *handler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extracting the id from the request url
	id := mux.Vars(r)["id"]

	var err error
	var user *models.User

	userid, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: id must be an integer")))

		return
	}

	user, err = uh.service.List(userid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	val, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(val)
}

// ListAll calls the listall function to fetch all the records from the database.
func (uh *handler) ListAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error

	var users []*models.User

	users, err = uh.service.ListAll()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	val, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(val)
}
