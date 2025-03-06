package postgres

import (
	// "AhmadAbdelrazik/arbun/internal/domain"
	"database/sql"
	"errors"
)

var (
	ErrDuplicateProduct          = errors.New("duplicate product")
	ErrProductNotFound           = errors.New("product not found")
	ErrEditConflict              = errors.New("edit conflict")
	ErrInsufficientProductAmount = errors.New("insufficient amount")
)

type ProductModel struct {
	*sql.DB
}

// func (m *ProductModel) InsertProduct(p domain.Product) (domain.Product, error) {
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return domain.Product{}, err
// 	}
// 	query := `
// 	INSERT INTO
// 	products(name, description, properties, price, amount)
// 	VALUES ($1, $2, $3, $4, $5)
// 	RETURNING (id, version)
// 	`
// }
