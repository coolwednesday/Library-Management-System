package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

var S *sql.DB

// User Structure
type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

// Add receives json object of user details and calls add function to add the record in the database
func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//b is the slice of bytes to read the request body
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	_ = json.Unmarshal(b, &user)
	err = user.add()
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte(fmt.Sprintf("Error: %v", err))))
		return

	} else {

		w.WriteHeader(http.StatusCreated)
		val := struct {
			Message string `json:"message"`
			Id      int    `json:"id"`
		}{"User added successfully", user.Id}
		v, _ := json.Marshal(val)
		log.Println(w.Write(v))
		return

	}

}

// Remove receives user id and calls remove function to remove the record from the database
func Remove(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//extracting the id from the request path
	id := mux.Vars(r)["id"]

	var err error
	var user User
	user.Id, err = strconv.Atoi(id)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte(fmt.Sprint("Error: ID must be an integer"))))
		return

	}

	val, err := user.remove()
	if err != nil {

		if val == "404" {

			w.WriteHeader(http.StatusNotFound)

		} else {

			w.WriteHeader(http.StatusBadRequest)

		}
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
	return

}

// List receives user id and calls list function to list the record from the database
func List(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//extracting the id from the request url
	id := mux.Vars(r)["id"]

	var err error

	user := &User{}
	user.Id, err = strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ID must be an integer")))
		return
	}

	user, err = user.List()

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return

	}

	val, err := json.Marshal(user)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return

	}

	w.WriteHeader(http.StatusOK)
	w.Write(val)

}

// ListAllUserHandler calls the listall function to fetch all the records from the database
func ListAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err error

	var users []*User

	user := User{}
	users, err = user.listall()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
	}
	val, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(val)

}

// add function
func (u *User) add() error {
	//Connecting to database
	/*
		s, err := database_conn.NewConnection()
		if err != nil {
			return err
		}
		defer s.Close()
	*/
	if u.Id == 0 {
		return fmt.Errorf("User Details Required")
	}
	if u.Id/1000 < 1 || u.Id/1000 >= 10 {
		return fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	_, err := S.Exec(`INSERT INTO users (name,Id) VALUES (?,?)`, u.Name, u.Id)
	if err != nil {
		return fmt.Errorf("Duplicate ID. User Already exist. Try again!")
	}
	return nil
}

// removeBook function
func (u *User) remove() (string, error) {
	/*
		s, err := database_conn.NewConnection()
		if err != nil {
			return "", err
		}
		defer s.Close()
	*/
	if u.Id/1000 < 1 || u.Id/1000 >= 10 {
		return "", fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	result, err := S.Exec(`DELETE FROM users WHERE id=?`, u.Id)
	if err != nil {
		return "", fmt.Errorf("User cannot be removed. User must return the book before being removed.")

	} else {
		val, err := result.RowsAffected()
		if err != nil {

			return "", err
		} else if val == 0 {
			return "404", fmt.Errorf("User with this ID does not exist.")
		}
	}

	return "User removed successfully", nil

}

// listBook function
func (u *User) List() (*User, error) {
	/*
		s, err := database_conn.NewConnection()
		if err != nil {
			return &User{}, err
		}
		defer s.Close()
	*/
	if u.Id/1000 < 1 || u.Id/1000 >= 10 {
		return &User{}, fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	err := S.QueryRow(`Select id,name from users where id=?`, u.Id).Scan(&u.Id, &u.Name)
	if err != nil {

		if err == sql.ErrNoRows {

			return &User{}, errors.New("User with this ID does not exist.")
		}
		return &User{}, err
	}
	return u, nil
}

// listall function
func (u *User) listall() ([]*User, error) {
	/*
		s, err := database_conn.NewConnection()
		if err != nil {
			return nil, err
		}
		defer s.Close()
	*/

	books := make([]*User, 0)
	rows, err := S.Query(`Select id,name from users where deleted_at is null;`)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*User{}, errors.New("No Users available.")
		}
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name)
		if err != nil {
			return []*User{}, err
		}
		books = append(books, u)
	}
	return books, nil
}
