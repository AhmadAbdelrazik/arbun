package models

import (
	"errors"
	"os/user"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrDuplicateUser = errors.New("duplicate user")
)

type UserModel struct {
	users     []user.User
	idCounter int64
}

func newUserModel() *UserModel {
	return &UserModel{
		users:     make([]user.User, 0),
		idCounter: 1,
	}
}

func (m *UserModel) InsertUser(c user.User) (user.User, error) {
	for _, a := range m.users {
		if a.Email == c.Email {
			return user.User{}, ErrDuplicateUser
		}
	}

	c.ID = m.idCounter
	c.Version = 1
	m.users = append(m.users, c)
	m.idCounter++

	return c, nil
}

func (m *UserModel) GetUserByEmail(email string) (user.User, error) {
	for _, a := range m.users {
		if a.Email == email {
			return a, nil
		}
	}

	return user.User{}, ErrUserNotFound
}

func (m *UserModel) GetUserByID(id int64) (user.User, error) {
	for _, a := range m.users {
		if a.ID == id {
			return a, nil
		}
	}

	return user.User{}, ErrUserNotFound
}

func (m *UserModel) GetAllUsers() ([]user.User, error) {
	return m.users, nil
}

func (m *UserModel) UpdateUser(c user.User) (user.User, error) {
	for i, cc := range m.users {
		if cc.ID == c.ID {
			if cc.Version != c.Version {
				return user.User{}, ErrEditConflict
			}

			c.Version = cc.Version + 1
			m.users[i] = c
			return c, nil
		}
	}

	return user.User{}, ErrUserNotFound
}

func (m *UserModel) DeleteUser(id int64) error {
	for i, a := range m.users {
		if a.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}
