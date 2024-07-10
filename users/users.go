package users

import "fmt"

type User struct {
	name string
	id   int
}

func AddUsers(existing []User, name string, id int) ([]User, error) {
	for _, user := range existing {
		if user.id == id {
			if user.name == name {
				return existing, fmt.Errorf("User already exists")
			}
			return existing, fmt.Errorf("User with %v id already exists.Try Again", id)
		}
	}
	existing = append(existing, User{name, id})
	return existing, nil
}

func RemoveUser(existing []User, id int) ([]User, error) {

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

func ListUser(existing []User, id int) (User, error) {

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
