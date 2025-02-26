package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"errors"
)

var (
	ErrAdminNotFound  = errors.New("admin not found")
	ErrDuplicateAdmin = errors.New("duplicate admin")
)

type AdminModel struct {
	admins    []domain.Admin
	idCounter int64
}

func newAdminModel() *AdminModel {
	return &AdminModel{
		admins:    make([]domain.Admin, 0),
		idCounter: 1,
	}
}

func (m *AdminModel) InsertAdmin(a domain.Admin) (domain.Admin, error) {
	for _, aa := range m.admins {
		if aa.Email == a.Email {
			return domain.Admin{}, ErrDuplicateAdmin
		}
	}

	a.ID = m.idCounter
	a.Version = 1
	m.admins = append(m.admins, a)
	m.idCounter++

	return a, nil
}

func (m *AdminModel) GetAdminByEmail(email string) (domain.Admin, error) {
	for _, aa := range m.admins {
		if aa.Email == email {
			return aa, nil
		}
	}

	return domain.Admin{}, ErrAdminNotFound
}

func (m *AdminModel) GetAdminByID(id int64) (domain.Admin, error) {
	for _, aa := range m.admins {
		if aa.ID == id {
			return aa, nil
		}
	}

	return domain.Admin{}, ErrAdminNotFound
}

func (m *AdminModel) GetAllAdmins() ([]domain.Admin, error) {
	return m.admins, nil
}

func (m *AdminModel) UpdateAdmin(a domain.Admin) (domain.Admin, error) {
	for i, aa := range m.admins {
		if aa.ID == a.ID {
			if aa.Version != a.Version {
				return domain.Admin{}, ErrEditConflict
			}

			a.Version = aa.Version + 1
			m.admins[i] = a
			return a, nil
		}
	}

	return domain.Admin{}, ErrAdminNotFound
}

func (m *AdminModel) DeleteAdmin(id int64) error {
	for i, aa := range m.admins {
		if aa.ID == id {
			m.admins = append(m.admins[:i], m.admins[i+1:]...)
			return nil
		}
	}
	return ErrAdminNotFound
}
