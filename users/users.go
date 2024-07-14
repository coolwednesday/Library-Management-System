package users

import (
	"SimpleRESTApi/database_conn"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"strconv"
)

// User Structure
type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

// AddUserHandler receives json object of user details and calls addUser function to add the record in the database
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	err = json.Unmarshal(b, &user)

	switch r.Method {
	case http.MethodPost:
		err = user.addUser()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("User added successfully"))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// RemoveUserHandler receives user id and calls removeUser function to remove the record from the database
func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	var user User
	user.Id, err = strconv.Atoi(string(b))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ID must be an integer")))
		return
	}

	switch r.Method {
	case http.MethodDelete:
		val, err := user.removeUser()
		if err != nil {
			if val == "404" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// ListUserHandler receives user id and calls listUser function to list the record from the database
func ListUserHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)

	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	user := &User{}
	user.Id, err = strconv.Atoi(string(b))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Error: ID must be an integer")))
		return
	}

	switch r.Method {
	case http.MethodGet:

		user, err = user.listUser()

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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}

}

// ListAllUserHandler calls the listAllUsers function to fetch all the records from the database
func ListAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
	}

	var users []*User
	switch r.Method {
	case http.MethodGet:
		user := User{}
		users, err = user.listAllUsers()
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}
}

// addUser function
func (u *User) addUser() error {
	//Connecting to database
	s, err := database_conn.NewConnection()
	if err != nil {
		return err
	}
	defer s.Close()

	if u.Id == 0 {
		return fmt.Errorf("User Details Required")
	}
	if u.Id/1000 < 1 || u.Id/1000 >= 10 {
		return fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	_, err = s.Exec(`INSERT INTO users (name,Id) VALUES (?,?)`, u.Name, u.Id)
	if err != nil {
		return fmt.Errorf("Duplicate ID. User Already exist. Try again!")
	}
	return nil
}

// removeBook function
func (u *User) removeUser() (string, error) {

	s, err := database_conn.NewConnection()
	if err != nil {
		return "", err
	}
	defer s.Close()

	if u.Id/1000 < 1 || u.Id/1000 >= 10 {
		return "", fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}

	result, err := s.Exec(`DELETE FROM users WHERE id=?`, u.Id)
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
func (u *User) listUser() (*User, error) {

	s, err := database_conn.NewConnection()
	if err != nil {
		return &User{}, err
	}
	defer s.Close()

	if u.Id/1000 < 1 || u.Id/1000 >= 10 {
		return &User{}, fmt.Errorf("Enter valid ID. Must be 4 Digits only.")
	}
	err = s.QueryRow(`Select * from users where id=?`, u.Id).Scan(&u.Id, &u.Name)
	if err != nil {

		if err == sql.ErrNoRows {

			return &User{}, errors.New("User with this ID does not exist.")
		}
		return &User{}, err
	}
	return u, nil
}

// listAllUsers function
func (u *User) listAllUsers() ([]*User, error) {
	s, err := database_conn.NewConnection()
	if err != nil {
		return nil, err
	}
	defer s.Close()
	books := make([]*User, 0)
	rows, err := s.Query(`Select id,name from users;`)
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
