package users

type UserStorer interface {
	add(int, string) error
	remove(int) (string, error)
	list(int) (*User, error)
	listall() ([]*User, error)
}
