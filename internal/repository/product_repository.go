package repository

import (
	"database/sql"

	"github.com/alvindashahrul/my-app/internal/model"
)

type ProductRepository interface {
	GetAllProducts() ([]model.Products, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAllProducts() ([]model.Products, error) {
	rows, err := r.db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Products
	for rows.Next() {
		var product model.Products
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
