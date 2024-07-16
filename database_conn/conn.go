package databaseconn

import (
	"database/sql"

	// driver to connect to sql
	_ "github.com/go-sql-driver/mysql"
)

func NewConnection() (*sql.DB, error) {
	// Connecting to the mysql Database
	s, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/library")
	if err != nil {
		return s, err
	}
	return s, nil
}
