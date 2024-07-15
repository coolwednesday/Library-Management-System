package main

import (
	"SimpleRESTApi/books"
	"SimpleRESTApi/users"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

/*
Mini Project Requirements:
Books - Add, List, ListAll
Users - Add, List, Remove, Borrow, Return
Use mysql database
REST APIS for Each Function
Write HTTP Tests for Them
*/

// userHandler function to route to the appropriate handler
func userHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		users.ListUserHandler(w, r)
	case http.MethodPost:
		users.AddUserHandler(w, r)
	case http.MethodDelete:
		users.RemoveUserHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}
}

// bookHandler function to route to the appropriate route
func bookHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		books.ListBookHandler(w, r)
	case http.MethodPost:
		books.AddBookHandler(w, r)
	case http.MethodPut:
		books.RemoveBookHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}
}

// bookRentHandler function
func bookRentHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		books.BorrowBookHandler(w, r)
	case http.MethodPost:
		books.ReturnBookHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("HTTP method %q not allowed", r.Method)))
	}
}

func main() {
	//routes with base url : http://localhost:8080/
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/book/all", books.ListAvailibleBooksHandler)
	http.HandleFunc("/user/all", users.ListAllUsersHandler)
	http.HandleFunc("/book/rent", bookRentHandler)

	//Connecting to Server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	} else {
		log.Println("Server successfully started. Listening on port 8080")
	}

}
