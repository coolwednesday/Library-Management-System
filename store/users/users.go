package users

import (
	"SimpleRESTApi/models"
	"database/sql"
	// connecting sql driver.
	_ "github.com/go-sql-driver/mysql"
)

// UserStore struct.
type store struct {
	DB *sql.DB
}

func New(db *sql.DB) *store {
	return &store{db}
}

func (s *store) CheckUser(id int) error {
	err := s.DB.QueryRow(`Select id from users where id=? and deleted_at is not null`, id).Scan(&id)

	return err
}

func (s *store) UpdateUser(id int) error {
	_, err := s.DB.Exec(`UPDATE users set deleted_at=null where id=?`, id)
	if err != nil {
		return err
	}
	return nil
}

// add function.
func (s *store) Add(id int, name string) error {
	_, err := s.DB.Exec(`INSERT INTO users (name,Id) VALUES (?,?)`, name, id)
	return err
}

func (s *store) Remove(id int) (sql.Result, error) {
	result, err := s.DB.Exec(`UPDATE users SET deleted_at = now() WHERE id = ? and deleted_at is null`, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// list function.
func (s *store) List(id int) (*models.User, error) {
	var name string

	err := s.DB.QueryRow(`Select name from users where id=? and deleted_at is null`, id).Scan(&name)
	if err != nil {
		return nil, err
	}

	return &models.User{
		name,
		id,
	}, nil
}

// listall function.
func (s *store) ListAll() ([]*models.User, error) {
	rows, err := s.DB.Query(`Select id,name from users where deleted_at is null;`)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Id, &u.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}
