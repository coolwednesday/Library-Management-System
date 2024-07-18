package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	// connecting sql driver.
	_ "github.com/go-sql-driver/mysql"
)

var S *sql.DB

// UserHandler defines a struct that has a UserStorer interface.
type UserHandler struct {
	UserStore UserStorer
}

// UserStore struct.
type UserStore struct {
}

// User Structure.
type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

// Add receives json object of user details and calls add function to add the record in the database.
func (uh *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// b is the slice of bytes to read the request body
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	_ = json.Unmarshal(b, &user)

	if user.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: user details required"))

		return
	}

	if user.Id/1000 < 1 || user.Id/1000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: enter valid id. Must be 4 digits only"))

		return
	}

	err = uh.UserStore.add(user.Id, user.Name)
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
func (uh *UserHandler) Remove(w http.ResponseWriter, r *http.Request) {
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

	if userid/1000 < 1 || userid/1000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error: enter valid id. Must be 4 digits only"))

		return
	}

	val, err := uh.UserStore.remove(userid)
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
func (uh *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extracting the id from the request url
	id := mux.Vars(r)["id"]

	var err error
	var user *User

	userid, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: id must be an integer")))

		return
	}

	if userid/1000 < 1 || userid/1000 >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprint("error: enter valid id. Must be 4 digits only")))

		return
	}

	user, err = uh.UserStore.list(userid)

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
func (uh *UserHandler) ListAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error

	var users []*User

	users, err = uh.UserStore.listall()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

// add function.
func (u *UserStore) add(id int, name string) error {

	_, err := S.Exec(`INSERT INTO users (name,Id) VALUES (?,?)`, name, id)
	if err != nil {
		return errors.New("duplicate id. User already exist. Try again")
	}

	return nil
}

// removeBook function.
func (u *UserStore) remove(id int) (string, error) {
	result, err := S.Exec(`DELETE FROM users WHERE id=?`, id)
	if err != nil {
		return "", errors.New("user cannot be removed. User must return the book before being removed")
	}

	val, err := result.RowsAffected()
	if err != nil {
		return "", err
	} else if val == 0 {
		return "404", errors.New("user with this id does not exist")
	}

	return "User removed successfully", nil
}

// list function.
func (u *UserStore) list(id int) (*User, error) {
	var name string

	err := S.QueryRow(`Select name from users where id=?`, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user with this id does not exist")
		}

		return nil, err
	}

	return &User{
		name,
		id,
	}, nil
}

// listall function.
func (us *UserStore) listall() ([]*User, error) {
	user := make([]*User, 0)

	rows, err := S.Query(`Select id,name from users where deleted_at is null;`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*User{}, errors.New("no users available")
		}

		return nil, err
	}

	u := &User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name)
		if err != nil {
			return nil, err
		}

		user = append(user, u)
	}

	return user, nil
}
