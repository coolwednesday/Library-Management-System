package main

import (
	"awesomeProject/books"
	"awesomeProject/users"
	"fmt"
	"sync"
)

// LendBooks implements the mutual exclusion in goroutines .
// It takes a WaitGroup, processId referring to action of Borrow or return,
// lending records as a map and existing list of Books
func LendBooks(processId int, wg *sync.WaitGroup, lend map[int]int, b2 []books.Book) {
	for {
		b := books.Book{}
		//Checks for which process needs to be executed : 1-> BorrowBook, 2->ReturnBook
		if processId == 1 {

			fmt.Println("BorrowBook can proceed")

			v, err := b.BorrowBook(b2, lend, 1234, 1234)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(v)
			}
			wg.Done()

			//break from the loop when the action is completed
			break

		} else if processId == 2 {

			fmt.Println("ReturnBook can proceed")

			v, err := b.ReturnBook(b2, lend, 1234)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(v, "Book Borrowed Successfully")
			}

			wg.Done()

			//break from the loop when the action is completed
			break
		}
	}
}

func main() {
	v := []users.User{} // Empty List of Users

	//Adding new users
	v1, err := users.AddUsersHelper(v, "Divya", 1234)
	if err != nil {
		fmt.Println(err)
	}
	_, err = users.AddUsersHelper(v1, "Jack", 1345)
	if err != nil {
		fmt.Println(err)
	}

	b := []books.Book{} // Empty List of Books

	//Adding new books to the list
	b1, err := books.AddBooksHelper(b, 1234, "Book1", "author1")
	if err != nil {
		fmt.Println(err)
	}
	b2, err := books.AddBooksHelper(b1, 1456, "Book2", "author4")
	if err != nil {
		fmt.Println(err)
	}

	//Creating a records of map with isbn of books as key and user of id as value
	lend := make(map[int]int)
	lend[1234] = 1234

	//waitGroup to prevent termination before complete execution of goroutines
	var wg sync.WaitGroup

	//add waiting for 2 routines
	wg.Add(2)

	//go routines to Borrow Book
	go LendBooks(1, &wg, lend, b2)

	//go routines to Lend Book
	go LendBooks(2, &wg, lend, b2)

	//Wait for all goroutines to complete
	wg.Wait()

}
