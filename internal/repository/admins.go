package repository

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int64
	Email    string
	Password Password
	FullName string
	Version  int
}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.hash = hash
	p.plaintext = &password
	return nil
}

func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

var (
	ErrAdminNotFound  = errors.New("admin not found")
	ErrDuplicateAdmin = errors.New("duplicate admin")
)

type AdminModel struct {
	admins    []Admin
	idCounter int64
}

func NewAdminModel() *AdminModel {
	return &AdminModel{
		admins:    make([]Admin, 0),
		idCounter: 1,
	}
}

func (m *AdminModel) InsertAdmin(admin Admin) (Admin, error) {
	for _, a := range m.admins {
		if a.Email == admin.Email {
			return Admin{}, ErrDuplicateAdmin
		}
	}

	admin.ID = m.idCounter
	admin.Version = 1
	m.admins = append(m.admins, admin)
	m.idCounter++

	return admin, nil
}

func (m *AdminModel) GetAdminByEmail(email string) (Admin, error) {
	for _, a := range m.admins {
		if a.Email == email {
			return a, nil
		}
	}

	return Admin{}, ErrAdminNotFound
}

func (m *AdminModel) GetAdminByID(id int64) (Admin, error) {
	for _, a := range m.admins {
		if a.ID == id {
			return a, nil
		}
	}

	return Admin{}, ErrAdminNotFound
}

func (m *AdminModel) GetAllAdmins() ([]Admin, error) {
	return m.admins, nil
}

func (m *AdminModel) UpdateAdmin(admin Admin) (Admin, error) {
	for i, a := range m.admins {
		if a.ID == admin.ID {
			if a.Version != admin.Version {
				return Admin{}, ErrEditConflict
			}

			admin.Version = a.Version + 1
			m.admins[i] = admin
			return admin, nil
		}
	}

	return Admin{}, ErrAdminNotFound
}

func (m *AdminModel) DeleteAdmin(id int64) error {
	for i, a := range m.admins {
		if a.ID == id {
			m.admins = append(m.admins[:i], m.admins[i+1:]...)
			return nil
		}
	}
	return ErrAdminNotFound
}
