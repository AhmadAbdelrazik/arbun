package models

import (
	"AhmadAbdelrazik/arbun/internal/domain/admin"
	"errors"
)

var (
	ErrAdminNotFound  = errors.New("admin not found")
	ErrDuplicateAdmin = errors.New("duplicate admin")
)

type AdminModel struct {
	admins    []admin.Admin
	idCounter int64
}

func NewAdminModel() *AdminModel {
	return &AdminModel{
		admins:    make([]admin.Admin, 0),
		idCounter: 1,
	}
}

func (m *AdminModel) InsertAdmin(a admin.Admin) (admin.Admin, error) {
	for _, aa := range m.admins {
		if aa.Email == a.Email {
			return admin.Admin{}, ErrDuplicateAdmin
		}
	}

	a.ID = m.idCounter
	a.Version = 1
	m.admins = append(m.admins, a)
	m.idCounter++

	return a, nil
}

func (m *AdminModel) GetAdminByEmail(email string) (admin.Admin, error) {
	for _, aa := range m.admins {
		if aa.Email == email {
			return aa, nil
		}
	}

	return admin.Admin{}, ErrAdminNotFound
}

func (m *AdminModel) GetAdminByID(id int64) (admin.Admin, error) {
	for _, aa := range m.admins {
		if aa.ID == id {
			return aa, nil
		}
	}

	return admin.Admin{}, ErrAdminNotFound
}

func (m *AdminModel) GetAllAdmins() ([]admin.Admin, error) {
	return m.admins, nil
}

func (m *AdminModel) UpdateAdmin(a admin.Admin) (admin.Admin, error) {
	for i, aa := range m.admins {
		if aa.ID == a.ID {
			if aa.Version != a.Version {
				return admin.Admin{}, ErrEditConflict
			}

			a.Version = aa.Version + 1
			m.admins[i] = a
			return a, nil
		}
	}

	return admin.Admin{}, ErrAdminNotFound
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
