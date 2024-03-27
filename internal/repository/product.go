package repository

import (
	"database/sql"

	"main.go/internal/model"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetProductById(id string) (*model.Products, error) {
	var p model.Products
	row := r.DB.QueryRow(`SELECT * FROM products WHERE id = $1`, id)
	if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Description); err != nil {
		return nil, err
	}

	return &p, nil
}
