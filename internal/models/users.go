package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"

	"errors"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrDuplicateUser = errors.New("duplicate user")
)

type IUserModel interface {
	InsertUser(c domain.User) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByID(id int64) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(c domain.User) (domain.User, error)
	DeleteUser(id int64) error
}

type UserModel struct {
	users     []domain.User
	idCounter int64
}

func newUserModel() *UserModel {
	return &UserModel{
		users:     make([]domain.User, 0),
		idCounter: 1,
	}
}

func (m *UserModel) InsertUser(c domain.User) (domain.User, error) {
	for _, a := range m.users {
		if a.Email == c.Email {
			return domain.User{}, ErrDuplicateUser
		}
	}

	c.ID = m.idCounter
	c.Version = 1
	m.users = append(m.users, c)
	m.idCounter++

	return c, nil
}

func (m *UserModel) GetUserByEmail(email string) (domain.User, error) {
	for _, a := range m.users {
		if a.Email == email {
			return a, nil
		}
	}

	return domain.User{}, ErrUserNotFound
}

func (m *UserModel) GetUserByID(id int64) (domain.User, error) {
	for _, a := range m.users {
		if a.ID == id {
			return a, nil
		}
	}

	return domain.User{}, ErrUserNotFound
}

func (m *UserModel) GetAllUsers() ([]domain.User, error) {
	return m.users, nil
}

func (m *UserModel) UpdateUser(c domain.User) (domain.User, error) {
	for i, cc := range m.users {
		if cc.ID == c.ID {
			if cc.Version != c.Version {
				return domain.User{}, ErrEditConflict
			}

			c.Version = cc.Version + 1
			m.users[i] = c
			return c, nil
		}
	}

	return domain.User{}, ErrUserNotFound
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
