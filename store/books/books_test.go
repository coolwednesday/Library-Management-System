package books

import (
	"SimpleRESTApi/models"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	s := New(db)

	tests := []struct {
		name        string
		isbn        int
		title       string
		author      string
		mockExpect  func()
		expectedErr error
	}{
		{
			name:   "Successful Add",
			isbn:   12345,
			title:  "Book Title",
			author: "Book Author",
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Book Title", "Book Author", 12345).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:   "Error on Add",
			isbn:   12345,
			title:  "Book Title",
			author: "Book Author",
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO books").
					WithArgs("Book Title", "Book Author", 12345).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			err := s.Add(tt.isbn, tt.title, tt.author)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	s := New(db)

	tests := []struct {
		name        string
		isbn        int
		mockExpect  func()
		expectedErr error
	}{
		{
			name: "Successful Remove",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectExec("UPDATE books SET deleted_at").
					WithArgs(12345).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Error on Remove",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectExec("UPDATE books SET deleted_at").
					WithArgs(12345).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			_, err := s.Remove(tt.isbn)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	s := New(db)

	tests := []struct {
		name         string
		isbn         int
		mockExpect   func()
		expectedBook *models.Book
		expectedErr  error
	}{
		{
			name: "Successful List",
			isbn: 12345,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"title", "author"}).
					AddRow("Book Title", "Book Author")
				mock.ExpectQuery(`^Select title,author from books where isbn=\? and deleted_at is null$`).
					WithArgs(12345).
					WillReturnRows(rows)
			},
			expectedBook: &models.Book{
				Title:  "Book Title",
				Author: "Book Author",
				Isbn:   12345,
			},
			expectedErr: nil,
		},
		{
			name: "Error on List",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectQuery(`^Select title,author from books where isbn=\? and deleted_at is null$`).
					WithArgs(12345).
					WillReturnError(sql.ErrConnDone)
			},
			expectedBook: nil,
			expectedErr:  sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			book, err := s.List(tt.isbn)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if book != nil && tt.expectedBook != nil {
				if *book != *tt.expectedBook {
					t.Errorf("expected book %v, got %v", *tt.expectedBook, *book)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestCheckBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	s := New(db)

	tests := []struct {
		name        string
		isbn        int
		mockExpect  func()
		expectedErr error
	}{
		{
			name: "Book Exists",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectQuery(`^Select isbn from books where isbn=\? and deleted_at is not null$`).
					WithArgs(12345).
					WillReturnRows(sqlmock.NewRows([]string{"isbn"}).AddRow(12345))
			},
			expectedErr: nil,
		},
		{
			name: "Book Does Not Exist",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectQuery(`^Select isbn from books where isbn=\? and deleted_at is not null$`).
					WithArgs(12345).
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			err := s.CheckBook(tt.isbn)
			if err != nil && err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestBorrow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	s := New(db)

	tests := []struct {
		name        string
		userID      int
		isbn        int
		mockExpect  func()
		expectedErr error
	}{
		{
			name:   "Successful Borrow",
			userID: 5678,
			isbn:   12345,
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO lendingRecords").
					WithArgs(5678, 12345).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:   "Error on Borrow",
			userID: 5678,
			isbn:   12345,
			mockExpect: func() {
				mock.ExpectExec("INSERT INTO lendingRecords").
					WithArgs(5678, 12345).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			err := s.Borrow(tt.userID, tt.isbn)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestReturnbook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	s := New(db)

	tests := []struct {
		name        string
		isbn        int
		mockExpect  func()
		expectedErr error
	}{
		{
			name: "Successful Return",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectExec("UPDATE lendingRecords SET deleted_at").
					WithArgs(12345).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Error on Return",
			isbn: 12345,
			mockExpect: func() {
				mock.ExpectExec("UPDATE lendingRecords SET deleted_at").
					WithArgs(12345).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			_, err := s.Returnbook(tt.isbn)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name        string
		isbn        int
		expectedErr error
		mockExpect  func()
	}{
		{
			name:        "Successful update",
			isbn:        1234567890,
			expectedErr: nil,
			mockExpect: func() {
				mock.ExpectExec(`^UPDATE books set deleted_at=null where isbn=\?$`).
					WithArgs(1234567890).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Mock result with one row affected
			},
		},
		{
			name:        "Update error",
			isbn:        9876543210,
			expectedErr: sql.ErrConnDone,
			mockExpect: func() {
				mock.ExpectExec(`^UPDATE books set deleted_at=null where isbn=\?$`).
					WithArgs(9876543210).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		err := store.UpdateBook(tt.isbn)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestCheckAvailibleBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name        string
		isbn        int
		expectedErr error
		mockExpect  func()
	}{
		{
			name:        "Book is available",
			isbn:        1234567890,
			expectedErr: nil,
			mockExpect: func() {
				mock.ExpectQuery(`^Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where s.isbn=\? and r.bookid is null or r.deleted_at is not null;$`).
					WithArgs(1234567890).
					WillReturnRows(sqlmock.NewRows([]string{"isbn"}).AddRow(1234567890))
			},
		},
		{
			name:        "Book is not available",
			isbn:        9876543210,
			expectedErr: sql.ErrNoRows,
			mockExpect: func() {
				mock.ExpectQuery(`^Select s.isbn from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where s.isbn=\? and r.bookid is null or r.deleted_at is not null;$`).
					WithArgs(9876543210).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		err := store.CheckAvailibleBook(tt.isbn)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestListAvailible(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name         string
		expectedRows []struct {
			isbn   int
			title  string
			author string
		}
		expectedErr error
		mockExpect  func()
	}{
		{
			name: "Successful retrieval",
			expectedRows: []struct {
				isbn   int
				title  string
				author string
			}{
				{isbn: 1234567890, title: "Book Title 1", author: "Author 1"},
				{isbn: 9876543210, title: "Book Title 2", author: "Author 2"},
			},
			expectedErr: nil,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"isbn", "title", "author"}).
					AddRow(1234567890, "Book Title 1", "Author 1").
					AddRow(9876543210, "Book Title 2", "Author 2")
				mock.ExpectQuery(`^Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where r.bookid is null and s.deleted_at is null;$`).
					WillReturnRows(rows)
			},
		},
		{
			name:         "No available books",
			expectedRows: nil,
			expectedErr:  sql.ErrNoRows,
			mockExpect: func() {
				mock.ExpectQuery(`^Select s.isbn,s.title,s.author from books s LEFT JOIN lendingRecords r on s.isbn = r.bookid where r.bookid is null and s.deleted_at is null;$`).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		rows, err := store.ListAvailible()

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)

			var actualRows []struct {
				isbn   int
				title  string
				author string
			}
			for rows.Next() {
				var isbn int
				var title, author string
				err := rows.Scan(&isbn, &title, &author)
				if err != nil {
					t.Errorf("%v: error scanning row: %v", tt.name, err)
					continue
				}
				actualRows = append(actualRows, struct {
					isbn   int
					title  string
					author string
				}{isbn, title, author})
			}
			assert.ElementsMatch(t, tt.expectedRows, actualRows, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}
