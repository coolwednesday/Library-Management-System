package books

import (
	"fmt"
)

/*
Books Management:
Each book has a title, author, and a unique ISBN.
Books can be added, removed, and listed.
*/
//Book Structure
type Book struct {
	title  string
	author string
	ISBN   int
}

// AddBooks Function takes the list of present Books and details of the new book to be added and returns new list of availible books and error message
func AddBooks(existing []Book, ISBN int, title, author string) ([]Book, error) {
	for _, book := range existing {
		if book.ISBN == ISBN {
			if book.title == title {
				if book.author == author {
					return existing, fmt.Errorf("Book already exists")
				}
			}
			return existing, fmt.Errorf("Book with %v ISBN already exists.Try Again", ISBN)
		}
	}
	existing = append(existing, Book{title, author, ISBN})
	return existing, nil
}

// RemoveBook Function takes list of present books and ISBN of the book t be removed and returns modified list of existing Books and error if any
func RemoveBook(existing []Book, ISBN int) ([]Book, error) {

	if len(existing) == 0 {
		return existing, fmt.Errorf("No Books are present")
	}
	for i, book := range existing {
		if book.ISBN == ISBN {
			existing = append(existing[:i], existing[i+1:]...)
			return existing, nil
		}
	}
	return existing, fmt.Errorf("Book Not Found")

}

// ListBook takes the existing Books and the title of book being searched and returns error
func ListBook(existing []Book, title string) ([]Book, error) {

	if len(existing) == 0 {
		return existing, fmt.Errorf("No Books are present")
	}
	k := 0
	var s = make([]Book, 0)
	for _, book := range existing {
		if book.title == title {
			s = append(s, book)
			k++
		}
	}
	if k > 0 {
		return s, nil
	}
	return []Book{}, fmt.Errorf("Book Not Found")

}

// BorrowBook takes the list of existing books and a map that maps ISBN of book and id of user and returns Book borrowed and error if any
func BorrowBook(existing []Book, m map[int]int, ISBN, id int) (Book, error) {

	if len(existing) == 0 {
		return Book{}, fmt.Errorf("No Books are present. Cannot Borrow")
	}

	_, ok := m[ISBN]
	if ok {
		return Book{}, fmt.Errorf("Book Already borrowed")
	}
	for _, book := range existing {
		if book.ISBN == ISBN {
			m[ISBN] = id
			return book, nil
		}
	}
	return Book{}, fmt.Errorf("Book Not Found")

}

// ReturnBook takes the list of availible Books and map that maps book ISBN to user id  and returns update as a string and error if any
func ReturnBook(existing []Book, m map[int]int, ISBN, id int) (string, error) {
	if len(existing) == 0 {
		return "", fmt.Errorf("No Books Registered. Cannot Return")
	}
	_, ok := m[ISBN]
	if ok {
		delete(m, ISBN)
		return "Book Returned Successfully", nil
	}
	return "", fmt.Errorf("Book Not Registered")
}
