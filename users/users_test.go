package users

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddUsers(t *testing.T) {
	type res struct {
		existing []User
		err      error
	}
	type req struct {
		user     User
		existing []User
	}
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

	for _, test := range tests {
		v, err := AddUsers(test.input.existing, test.input.user.name, test.input.user.id)
		s := res{v, err}
		assert.Equal(t, test.output, s)
	}
}

func TestRemoveUser(t *testing.T) {
	type res struct {
		existing []User
		err      error
	}
	type req struct {
		id       int
		existing []User
	}
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

	for _, test := range tests {
		v, err := RemoveUser(test.input.existing, test.input.id)
		s := res{v, err}
		assert.Equal(t, test.output, s)
	}
}

func TestListUser(t *testing.T) {
	type res struct {
		found User
		err   error
	}
	type req struct {
		id       int
		existing []User
	}
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

	for _, test := range tests {
		v, err := ListUser(test.input.existing, test.input.id)
		s := res{v, err}
		assert.Equal(t, test.output, s)
	}
}
