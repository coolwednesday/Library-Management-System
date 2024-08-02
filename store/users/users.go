package users

import (
	"database/sql"
	"github.com/libraryManagementSystem/models"
	"gofr.dev/pkg/gofr"

	// connecting sql driver.
	_ "github.com/go-sql-driver/mysql"
)

// UserStore struct.
type store struct {
}

func New() *store {
	return &store{}
}

func (s *store) CheckUser(c *gofr.Context, id int) error {
	err := c.SQL.QueryRowContext(c, `Select id from users where id=? and deleted_at is not null`, id).Scan(&id)

	return err
}

func (s *store) UpdateUser(c *gofr.Context, id int) error {
	_, err := c.SQL.ExecContext(c, `UPDATE users set deleted_at=null where id=?`, id)
	if err != nil {
		return err
	}
	return nil
}

// add function.
func (s *store) Add(c *gofr.Context, id int, name string) error {
	_, err := c.SQL.ExecContext(c, `INSERT INTO users (name,Id) VALUES (?,?)`, name, id)
	return err
}

func (s *store) Remove(c *gofr.Context, id int) (sql.Result, error) {
	result, err := c.SQL.ExecContext(c, `UPDATE users SET deleted_at = now() WHERE id = ? and deleted_at is null`, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// list function.
func (s *store) List(c *gofr.Context, id int) (*models.User, error) {
	var name string

	err := c.SQL.QueryRowContext(c, `Select name from users where id=? and deleted_at is null`, id).Scan(&name)
	if err != nil {
		return nil, err
	}

	return &models.User{
		name,
		id,
	}, nil
}

// listall function.
func (s *store) ListAll(c *gofr.Context) ([]*models.User, error) {
	rows, err := c.SQL.QueryContext(c, `Select id,name from users where deleted_at is null;`)
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
