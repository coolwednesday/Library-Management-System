package main

import (
	bookHandler "github.com/libraryManagementSystem/handler/books"
	userHandler "github.com/libraryManagementSystem/handler/users"
	bookService "github.com/libraryManagementSystem/service/books"
	userService "github.com/libraryManagementSystem/service/users"
	bookStore "github.com/libraryManagementSystem/store/books"
	userStore "github.com/libraryManagementSystem/store/users"
	"gofr.dev/pkg/gofr"
)

func main() {
	// initialise gofr object
	app := gofr.New()

	bs := bookStore.New()
	bsvc := bookService.New(bs)
	bh := bookHandler.New(bsvc)

	// routes with base url : http://localhost:8080
	app.POST("/book", bh.Add)
	app.DELETE("/book/{isbn}", bh.Remove)
	app.GET("/book", bh.ListAvailible)
	app.GET("/book/{isbn}", bh.List)
	app.POST("/book/rent", bh.Borrow)
	app.DELETE("/book/rent/{isbn}", bh.Return)

	us := userStore.New()
	usvc := userService.New(us)
	uh := userHandler.New(usvc)

	app.POST("/user", uh.Add)
	app.GET("/user/{id}", uh.List)
	app.GET("/user", uh.ListAll)
	app.DELETE("/user/{id}", uh.Remove)

	app.Run()
}
