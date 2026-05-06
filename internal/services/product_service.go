package services

import (
	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
)

type ProductService interface {
	GetAllProducts() ([]model.Products, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts() ([]model.Products, error) {
	return s.repo.GetAllProducts()
}
