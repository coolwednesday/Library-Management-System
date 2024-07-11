package users

import "fmt"

// User struct
type User struct {
	name string
	id   int
}

// AddUser method add the User into existing list of users and return the modifies list and error if any
func (u User) AddUsers(existing []User) ([]User, error) {
	for _, user := range existing {
		if user.id == u.id {
			if user.name == u.name {
				return existing, fmt.Errorf("User already exists")
			}
			return existing, fmt.Errorf("User with %v id already exists.Try Again", u.id)
		}
	}
	existing = append(existing, u)
	return existing, nil
}

// RemoveUser method removes the user by their user id and then returns modified list of users and errors if any
func (u User) RemoveUser(existing []User, id int) ([]User, error) {

	if len(existing) == 0 {
		return existing, fmt.Errorf("No Users are present")
	}
	for i, user := range existing {
		if user.id == id {
			existing = append(existing[:i], existing[i+1:]...)
			return existing, nil
		}
	}
	return existing, fmt.Errorf("User Not Found")

}

// ListUser method returns the user that has the same id as argument as well as errors if any
func (u User) ListUser(existing []User, id int) (User, error) {

	if len(existing) == 0 {
		return User{}, fmt.Errorf("No Users are present")
	}

	for _, user := range existing {
		if user.id == id {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("User Not Found")

}
