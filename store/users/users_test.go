package users

import (
	"SimpleRESTApi/models"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestListAll(t *testing.T) {
	// making connection to mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}

	// closing the connection.
	defer db.Close()

	// initializing the store with mock db.
	store := New(db)

	// defining test cases.
	tests := []struct {
		name         string
		expectedRows []*models.User
		expectedErr  error
		mockExpect   func()
	}{
		{
			name: "Success with multiple users",
			expectedRows: []*models.User{
				{"User1", 1234},
				{"User2", 5678},
			},
			expectedErr: nil,
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1234, "User1").
					AddRow(5678, "User2")
				mock.ExpectQuery(`^Select id,name from users where deleted_at is null;$`).
					WillReturnRows(rows)
			},
		},
		{
			name:         "Error during query",
			expectedRows: nil,
			expectedErr:  sql.ErrConnDone,
			mockExpect: func() {
				mock.ExpectQuery(`^Select id,name from users where deleted_at is null;$`).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		var users []*models.User

		users, err := store.ListAll()

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.ElementsMatch(t, tt.expectedRows, users, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestCheckUser(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}

	defer db.Close()

	store := New(db)

	tests := []struct {
		name        string
		id          int
		expectedErr error
		mockExpect  func()
	}{
		{
			name:        "User exists and is marked as deleted",
			id:          1234,
			expectedErr: nil,
			mockExpect: func() {
				mock.ExpectQuery(`^Select id from users where id=\? and deleted_at is not null$`).
					WithArgs(1234).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1234))
			},
		},
		{
			name:        "User does not exist",
			id:          5678,
			expectedErr: sql.ErrNoRows,
			mockExpect: func() {
				mock.ExpectQuery(`^Select id from users where id=\? and deleted_at is not null$`).
					WithArgs(5678).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "Query error",
			id:          9999,
			expectedErr: sql.ErrConnDone,
			mockExpect: func() {
				mock.ExpectQuery(`^Select id from users where id=\? and deleted_at is not null$`).
					WithArgs(9999).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		err := store.CheckUser(tt.id)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name        string
		id          int
		nameStr     string
		expectedErr error
		mockExpect  func()
	}{
		{
			name:        "Successful insertion",
			id:          1234,
			nameStr:     "User1",
			expectedErr: nil,
			mockExpect: func() {
				mock.ExpectExec(`^INSERT INTO users \(name,Id\) VALUES \(\?,\?\)$`).
					WithArgs("User1", 1234).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Mock result with one row affected
			},
		},
		{
			name:        "Insertion error",
			id:          5678,
			nameStr:     "User2",
			expectedErr: sql.ErrConnDone,
			mockExpect: func() {
				mock.ExpectExec(`^INSERT INTO users \(name,Id\) VALUES \(\?,\?\)$`).
					WithArgs("User2", 5678).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		err := store.Add(tt.id, tt.nameStr)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestRemove(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name        string
		id          int
		expectedErr error
		mockExpect  func()
	}{
		{
			name:        "Successful removal",
			id:          1234,
			expectedErr: nil,
			mockExpect: func() {
				mock.ExpectExec(`^UPDATE users SET deleted_at = now\(\) WHERE id = \? and deleted_at is null$`).
					WithArgs(1234).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Mock result with one row affected
			},
		},
		{
			name:        "Removal error",
			id:          5678,
			expectedErr: sql.ErrConnDone,
			mockExpect: func() {
				mock.ExpectExec(`^UPDATE users SET deleted_at = now\(\) WHERE id = \? and deleted_at is null$`).
					WithArgs(5678).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		result, err := store.Remove(tt.id)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.NotNil(t, result, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name         string
		id           int
		expectedUser *models.User
		expectedErr  error
		mockExpect   func()
	}{
		{
			name:         "Successful retrieval",
			id:           3456,
			expectedUser: &models.User{"User5", 3456},
			expectedErr:  nil,
			mockExpect: func() {
				mock.ExpectQuery(`^Select name from users where id=\? and deleted_at is null$`).
					WithArgs(3456).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User5"))
			},
		},
		{
			name:         "User not found",
			id:           1234,
			expectedUser: nil,
			expectedErr:  sql.ErrNoRows,
			mockExpect: func() {
				mock.ExpectQuery(`^Select name from users where id=\? and deleted_at is null$`).
					WithArgs(1234).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		user, err := store.List(tt.id)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.Equal(t, tt.expectedUser, user, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	store := New(db)

	tests := []struct {
		name        string
		id          int
		expectedErr error
		mockExpect  func()
	}{
		{
			name:        "Successful update",
			id:          1234,
			expectedErr: nil,
			mockExpect: func() {
				mock.ExpectExec(`^UPDATE users set deleted_at=null where id=\?$`).
					WithArgs(1234).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Mock result with one row affected
			},
		},
		{
			name:        "Update error",
			id:          5678,
			expectedErr: sql.ErrConnDone,
			mockExpect: func() {
				mock.ExpectExec(`^UPDATE users set deleted_at=null where id=\?$`).
					WithArgs(5678).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt.mockExpect()

		err := store.UpdateUser(tt.id)

		if tt.expectedErr != nil {
			assert.ErrorIs(t, err, tt.expectedErr, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}

		assert.NoError(t, mock.ExpectationsWereMet(), tt.name+" - unmet expectations")
	}
}
