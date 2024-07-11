package books

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAddBooks test for all inputs in form of req struct and returns response in form of res struct
func TestAddBooks(t *testing.T) {
	//res struct defines updated existing slice of Books and the corresponding error if any
	type res struct {
		existing []Book
		err      error
	}
	//req struct defines initial existing slice of Books and the book details we want to add
	type req struct {
		book     Book
		existing []Book
	}
	//test cases for various usecases of Adding Books
	tests := []struct {
		name   string
		input  req
		output res
	}{
		{name: "No book is present",
			input:  req{Book{title: "Book1", author: "Author1", isbn: 1234}, []Book{}},
			output: res{[]Book{{title: "Book1", author: "Author1", isbn: 1234}}, nil},
		},

		{name: "Adding Book with same title",
			input: req{Book{title: "Book1", author: "Author2", isbn: 1357}, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}},
			output: res{[]Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}, errors.New("Book already exists")},
		},

		{name: "Adding Book with duplicate isbn",
			input: req{Book{title: "Book3", author: "Author1", isbn: 1234}, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}},
			output: res{[]Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}, errors.New("Book with 1234 isbn already exists.Try Again")},
		},
	}

	//Running for all testcases
	for _, test := range tests {

		v, err := AddBooksHelper(test.input.existing, test.input.book.isbn, test.input.book.title, test.input.book.author)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)

	}
}

// TestRemoveBook function tests for all inputs in form of req type and res type as output
func TestRemoveBook(t *testing.T) {
	//res struct describes the output of the updated slice of Books after removing book and the error if any
	type res struct {
		existing []Book
		err      error
	}
	//req describes the initial list of Books and the detail of the book that is to be removed
	type req struct {
		isbn     int
		existing []Book
	}
	//testcases for RemoveBook function
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{1234, []Book{}},
		output: res{[]Book{}, errors.New("No Books are present")},
	},
		{name: "Removing new different book",
			input: req{1234, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456}}},
			output: res{[]Book{
				{title: "Book2", author: "Author2", isbn: 1456}}, nil},
		},
		{name: "Book is not present",
			input: req{1567, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}},
			output: res{[]Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}, errors.New("Book Not Found")},
		},
	}
	//Creating an empty instance of Book struct to call the method
	b := Book{}
	//Testing all the tests for RemoveBook Function
	for _, test := range tests {
		v, err := b.RemoveBook(test.input.existing, test.input.isbn)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}

// TestListBook function tests for all inputs in form of req type and res type as output
func TestListBook(t *testing.T) {
	//res struct describes the list of the book with the title that we want to list as well as error if any
	type res struct {
		found []Book
		err   error
	}
	//req struct describes the list of the books and title of the book we wish to list
	type req struct {
		title    string
		existing []Book
	}
	//test cases for the ListBook Function
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{"Book1", []Book{}},
		output: res{[]Book{}, errors.New("No Books are present")},
	},
		{name: "Listing book",
			input: req{"Book1", []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456}}},
			output: res{[]Book{{title: "Book1", author: "Author1", isbn: 1234}}, nil},
		},
		{name: "Listing book with same names",
			input: req{"Book1", []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author3", isbn: 1900}}},
			output: res{[]Book{
				{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book1", author: "Author3", isbn: 1900},
			}, nil},
		},
		{name: "Listing a Book not present",
			input: req{"Book4", []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}},
			output: res{[]Book{}, errors.New("Book Not Found")},
		},
	}
	//creating an empty struct of Book type
	b := Book{}
	//testing all the testcases for ListBook function
	for _, test := range tests {
		v, err := b.ListBook(test.input.existing, test.input.title)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}

// TestBorrowBook function tests for all inputs in form of req type and res type as output
func TestBorrowBook(t *testing.T) {
	//Creating a map that maps has isbn of the book as key and id of user as value
	lend := make(map[int]int)
	//Adding a record
	lend[1234] = 1234
	//res struct has the book borrowed and the error if any
	type res struct {
		found Book
		err   error
	}
	//req struct describes the input that has id of user, isbn of book
	//as well as existing list of books and the current lending records as map
	type req struct {
		isbn     int
		id       int
		m        map[int]int
		existing []Book
	}
	//testcases for BorrowBook Function
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{1234, 1234, lend, []Book{}},
		output: res{Book{}, errors.New("No Books are present. Cannot Borrow")},
	},
		{name: "Borrowing new book",
			input: req{1456, 1234, lend, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456}}},
			output: res{Book{title: "Book2", author: "Author2", isbn: 1456}, nil},
		},
		{name: "Borrowing same book",
			input: req{1234, 1234, lend, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456}}},
			output: res{Book{}, errors.New("Book Already borrowed")},
		},
		{name: "Book is not present",
			input: req{1890, 1234, lend, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}},
			output: res{Book{}, errors.New("Book Not Found")},
		},
	}
	//Creating an empty struct of Book type
	b := Book{}
	//Checking for all testcases
	for _, test := range tests {
		v, err := b.BorrowBook(test.input.existing, test.input.m, test.input.id, test.input.isbn)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}

// TestReturnBook function tests for all inputs in form of req type and res type as output
func TestReturnBook(t *testing.T) {
	//Creating a map that maps has isbn of the book as key and id of user as value
	lend := make(map[int]int)
	//Adding a record
	lend[1234] = 1234
	//res struct describes a status update using returned string and the error if any
	type res struct {
		returned string
		err      error
	}
	//req struct describes the isbn of the book to be returned ,
	//existing slice of books and the current lrnding records as map
	type req struct {
		isbn     int
		m        map[int]int
		existing []Book
	}
	//testcases for ReturnBook function
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{1234, lend, []Book{}},
		output: res{"", errors.New("No Books Registered. Cannot Return")},
	},
		{name: "Returning Book",
			input: req{1234, lend, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456}}},
			output: res{"Book Returned Successfully", nil},
		},
		{name: "Book is not present in existing books",
			input: req{1890, lend, []Book{{title: "Book1", author: "Author1", isbn: 1234},
				{title: "Book2", author: "Author2", isbn: 1456},
				{title: "Book1", author: "Author2", isbn: 1357}}},
			output: res{"", errors.New("Book Not Registered")},
		},
	}
	//Creating an empty struct of Book type
	b := Book{}
	//Testing for all testcases
	for _, test := range tests {
		v, err := b.ReturnBook(test.input.existing, test.input.m, test.input.isbn)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}
