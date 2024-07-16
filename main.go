package main

import (
	"SimpleRESTApi/books"
	"database/sql"
	"github.com/gorilla/mux"
	//db "SimpleRESTApi/database_conn"
	"SimpleRESTApi/users"
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

func main() {
	r := mux.NewRouter()

	// creating UserHandler.
	uh := users.UserHandler{
		UserStore: &users.UserStore{},
	}

	// routes with base url : http://localhost:8080
	r.HandleFunc("/user", uh.Add).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", uh.List).Methods(http.MethodGet)
	r.HandleFunc("/user", uh.ListAll).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", uh.Remove).Methods(http.MethodDelete)

	// creating BookHandler.
	bh := books.BookHandler{
		BookStore: &books.BookStore{},
	}

	// routes with base url : http://localhost:8080
	r.HandleFunc("/book", bh.Add).Methods(http.MethodPost)
	r.HandleFunc("/book/{isbn}", bh.Remove).Methods(http.MethodDelete)
	r.HandleFunc("/book", bh.ListAvailible).Methods(http.MethodGet)
	r.HandleFunc("/book/{isbn}", bh.List).Methods(http.MethodGet)

	// creating RecordHandler.
	rh := books.RecordHandler{
		RentStore: &books.RentStore{},
	}

	// routes with base url : http://localhost:8080
	r.HandleFunc("/book/rent", rh.Borrow).Methods(http.MethodPost)
	r.HandleFunc("/book/rent/{isbn}", rh.Return).Methods(http.MethodDelete)
	http.Handle("/", r)

	var err error

	S, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/library")
	if err != nil {
		log.Println(err)
	} else {
		users.S = S
		books.S = S

		log.Println("Database connected")
	}

	//Connecting to Server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	} else {
		log.Println("Server successfully started. Listening on port 8080")
	}

}
