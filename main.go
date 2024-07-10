package main

import (
	"awesomeProject/books"
	"awesomeProject/users"
	"fmt"
	"sync"
	"time"
)

func main() {
	v := []users.User{} // Empty List of Users
	v1, err := users.AddUsers(v, "Divya", 1234)
	if err != nil {
		fmt.Println(err)
	}
	_, err = users.AddUsers(v1, "Jack", 1345)
	if err != nil {
		fmt.Println(err)
	}

	b := []books.Book{} // Empty List of Books
	b1, err := books.AddBooks(b, 1234, "Book1", "author1")
	if err != nil {
		fmt.Println(err)
	}
	b2, err := books.AddBooks(b1, 1456, "Book2", "author4")
	if err != nil {
		fmt.Println(err)
	}

	lend := make(map[int]int)
	lend[1234] = 1234

	ch := make(chan struct{}, 1) // Unbuffered channel for synchronization
	var wg sync.WaitGroup
	wg.Add(2) // Increase the WaitGroup count for two goroutines
	// Start the synchronization by sending a signal to the channel
	ch <- struct{}{}
	go func() {
		time.Sleep(2 * time.Second)
		defer wg.Done()
		<-ch // Wait for signal to proceed
		fmt.Println("BorrowBook can be processed")
		v, err := books.BorrowBook(b2, lend, 1456, 1234)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(v, "BorrowBook successfully processed")
		}
		// Signal completion to allow ReturnBook to proceed
		ch <- struct{}{}
	}()

	go func() {
		defer wg.Done()
		<-ch // Wait for BorrowBook to complete
		fmt.Println("ReturnBook can be processed")
		v, err := books.ReturnBook(b2, lend, 1234, 1234)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(v, "ReturnBook successfully processed")
		}
		ch <- struct{}{}
	}()

	wg.Wait()
	close(ch) // Close the channel once synchronization is done (optional, but recommended)
}
