package users

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAddUsers function tests for all inputs in form of req type and res type as output
func TestAddUsers(t *testing.T) {

	//res struct defines output that has modifies slice of users as well as error if any
	type res struct {
		existing []User
		err      error
	}

	//req struct defines input that has user as well as exiting slice of Users
	type req struct {
		user     User
		existing []User
	}

	//testcases for AddUserHelper function
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No User is present",
		input:  req{User{name: "User1", id: 1234}, []User{}},
		output: res{[]User{{name: "User1", id: 1234}}, nil},
	},
		{name: "Adding new User",
			input: req{User{name: "User2", id: 1345}, []User{
				{name: "User1", id: 1234},
			}},
			output: res{[]User{{name: "User1", id: 1234},
				{name: "User2", id: 1345}}, nil},
		},
		{name: "Adding same User",
			input: req{User{name: "User2", id: 1345}, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
			}},
			output: res{[]User{{name: "User1", id: 1234},
				{name: "User2", id: 1345}}, errors.New("User already exists")},
		},
		{name: "Adding User with same name",
			input: req{User{name: "User1", id: 1567}, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
			}},
			output: res{[]User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}, nil},
		},
		{name: "Adding User with duplicate id",
			input: req{User{name: "User4", id: 1567}, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}},
			output: res{[]User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}, errors.New("User with 1567 id already exists.Try Again")},
		},
	}

	//checking for all test cases
	for _, test := range tests {
		v, err := AddUsersHelper(test.input.existing, test.input.user.name, test.input.user.id)
		s := res{v, err}
		assert.Equal(t, test.output, s)
	}
}

// TestRemoveUser function tests for all inputs in form of req type and res type as output
func TestRemoveUser(t *testing.T) {

	//res struct defines the output that has existing users and error if any
	type res struct {
		existing []User
		err      error
	}

	//req struct defines the input that has the id of user and existing users
	type req struct {
		id       int
		existing []User
	}

	//testcases for RemoveUser method
	tests := []struct {
		name   string
		input  req
		output res
	}{
		{name: "No User is present",
			input:  req{1234, []User{}},
			output: res{[]User{}, errors.New("No Users are present")},
		},
		{name: "Removing User",
			input: req{1234, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
			}},
			output: res{[]User{
				{name: "User2", id: 1345},
			}, nil},
		},
		{name: "User is not present",
			input: req{1900, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}},
			output: res{[]User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}, errors.New("User Not Found")},
		},
	}

	//creating empty struct for User type
	u := User{}

	//testing for all testcases
	for _, test := range tests {
		v, err := u.RemoveUser(test.input.existing, test.input.id)
		s := res{v, err}
		assert.Equal(t, test.output, s)
	}
}

// TestListUser function tests for all inputs in form of req type and res type as output
func TestListUser(t *testing.T) {

	//res struct type defines the response with found user and error if any
	type res struct {
		found User
		err   error
	}

	//req struct type defines the id of user as well as the existing slice of User type
	type req struct {
		id       int
		existing []User
	}

	//testcases for ListUser method
	tests := []struct {
		name   string
		input  req
		output res
	}{{name: "No Users are present",
		input:  req{1345, []User{}},
		output: res{User{}, errors.New("No Users are present")},
	},
		{name: "Listing a User from id",
			input: req{1234, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}},
			output: res{User{name: "User1", id: 1234}, nil},
		},
		{name: "User is not present",
			input: req{1679, []User{
				{name: "User1", id: 1234},
				{name: "User2", id: 1345},
				{name: "User1", id: 1567},
			}},
			output: res{User{}, errors.New("User Not Found")},
		},
	}

	//creating an empty struct of User type
	u := User{}

	//testing for all test cases
	for _, test := range tests {
		v, err := u.ListUser(test.input.existing, test.input.id)
		s := res{v, err}
		assert.Equal(t, test.output, s)
	}
}
