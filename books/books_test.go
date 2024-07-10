package books

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddBooks(t *testing.T) {
	type res struct {
		existing []Book
		err      error
	}
	type req struct {
		book     Book
		existing []Book
	}
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{Book{title: "Book1", author: "Author1", ISBN: 1234}, []Book{}},
		output: res{[]Book{{title: "Book1", author: "Author1", ISBN: 1234}}, nil},
	},
		{name: "Adding Book with same title",
			input: req{Book{title: "Book1", author: "Author2", ISBN: 1357}, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}},
			output: res{[]Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}, errors.New("Book already exists")},
		},
		{name: "Adding Book with duplicate ISBN",
			input: req{Book{title: "Book3", author: "Author1", ISBN: 1234}, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}},
			output: res{[]Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}, errors.New("Book with 1234 ISBN already exists.Try Again")},
		},
	}

	for _, test := range tests {
		v, err := AddBooks(test.input.existing, test.input.book.ISBN, test.input.book.title, test.input.book.author)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)

	}
}

func TestRemoveBook(t *testing.T) {
	type res struct {
		existing []Book
		err      error
	}
	type req struct {
		ISBN     int
		existing []Book
	}
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{1234, []Book{}},
		output: res{[]Book{}, errors.New("No Books are present")},
	},
		{name: "Removing new different book",
			input: req{1234, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456}}},
			output: res{[]Book{
				{title: "Book2", author: "Author2", ISBN: 1456}}, nil},
		},
		{name: "Book is not present",
			input: req{1567, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}},
			output: res{[]Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}, errors.New("Book Not Found")},
		},
	}

	for _, test := range tests {
		v, err := RemoveBook(test.input.existing, test.input.ISBN)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}

func TestListBook(t *testing.T) {
	type res struct {
		found []Book
		err   error
	}
	type req struct {
		title    string
		existing []Book
	}
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{"Book1", []Book{}},
		output: res{[]Book{}, errors.New("No Books are present")},
	},
		{name: "Listing book",
			input: req{"Book1", []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456}}},
			output: res{[]Book{{title: "Book1", author: "Author1", ISBN: 1234}}, nil},
		},
		{name: "Listing book with same names",
			input: req{"Book1", []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author3", ISBN: 1900}}},
			output: res{[]Book{
				{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book1", author: "Author3", ISBN: 1900},
			}, nil},
		},
		{name: "Listing a Book not present",
			input: req{"Book4", []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}},
			output: res{[]Book{}, errors.New("Book Not Found")},
		},
	}

	for _, test := range tests {
		v, err := ListBook(test.input.existing, test.input.title)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}

func TestBorrowBook(t *testing.T) {
	lend := make(map[int]int)
	lend[1234] = 1234
	type res struct {
		found Book
		err   error
	}
	type req struct {
		ISBN     int
		id       int
		m        map[int]int
		existing []Book
	}
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{1234, 1234, lend, []Book{}},
		output: res{Book{}, errors.New("No Books are present. Cannot Borrow")},
	},
		{name: "Borrowing new book",
			input: req{1456, 1234, lend, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456}}},
			output: res{Book{title: "Book2", author: "Author2", ISBN: 1456}, nil},
		},
		{name: "Borrowing same book",
			input: req{1234, 1234, lend, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456}}},
			output: res{Book{}, errors.New("Book Already borrowed")},
		},
		{name: "Book is not present",
			input: req{1890, 1234, lend, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}},
			output: res{Book{}, errors.New("Book Not Found")},
		},
	}

	for _, test := range tests {
		v, err := BorrowBook(test.input.existing, test.input.m, test.input.ISBN, test.input.id)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}

func TestReturnBook(t *testing.T) {

	lend := make(map[int]int)
	lend[1234] = 1234

	type res struct {
		returned string
		err      error
	}

	type req struct {
		ISBN     int
		id       int
		m        map[int]int
		existing []Book
	}

	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No book is present",
		input:  req{1234, 1234, lend, []Book{}},
		output: res{"", errors.New("No Books Registered. Cannot Return")},
	},
		{name: "Returning Book",
			input: req{1234, 1234, lend, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456}}},
			output: res{"Book Returned Successfully", nil},
		},
		{name: "Book is not present in existing books",
			input: req{1890, 1234, lend, []Book{{title: "Book1", author: "Author1", ISBN: 1234},
				{title: "Book2", author: "Author2", ISBN: 1456},
				{title: "Book1", author: "Author2", ISBN: 1357}}},
			output: res{"", errors.New("Book Not Registered")},
		},
	}

	for _, test := range tests {
		v, err := ReturnBook(test.input.existing, test.input.m, test.input.ISBN, test.input.id)
		s := res{
			v, err,
		}
		assert.Equal(t, test.output, s)
	}
}
