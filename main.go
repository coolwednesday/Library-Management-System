package main

import (
	//"SimpleRESTApi/books"
	bookHandler "SimpleRESTApi/handler/books"
	userHandler "SimpleRESTApi/handler/users"
	bookService "SimpleRESTApi/services/books"
	userService "SimpleRESTApi/services/users"
	bookStore "SimpleRESTApi/store/books"
	userStore "SimpleRESTApi/store/users"
	"database/sql"
	"github.com/gorilla/mux"

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
	var err error
	r := mux.NewRouter()

	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/library")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database connected")
	}

	us := userStore.New(db)
	usvc := userService.New(us)
	uh := userHandler.New(usvc)

	bs := bookStore.New(db)
	bsvc := bookService.New(bs)
	bh := bookHandler.New(bsvc)

	// routes with base url : http://localhost:8080
	r.HandleFunc("/user", uh.Add).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", uh.List).Methods(http.MethodGet)
	r.HandleFunc("/user", uh.ListAll).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", uh.Remove).Methods(http.MethodDelete)

	// routes with base url : http://localhost:8080
	r.HandleFunc("/book", bh.Add).Methods(http.MethodPost)
	r.HandleFunc("/book/{isbn}", bh.Remove).Methods(http.MethodDelete)
	r.HandleFunc("/book", bh.ListAvailible).Methods(http.MethodGet)
	r.HandleFunc("/book/{isbn}", bh.List).Methods(http.MethodGet)
	r.HandleFunc("/book/rent", bh.Borrow).Methods(http.MethodPost)
	r.HandleFunc("/book/rent/{isbn}", bh.Return).Methods(http.MethodDelete)

	http.Handle("/", r)

	//Connecting to Server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	} else {
		log.Println("Server successfully started. Listening on port 8080")
	}

}
