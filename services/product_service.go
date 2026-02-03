package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"log/slog"
)

//intermediate lah disini sama repo dengan handler

type ProductService struct {
	repo   *repositories.ProductRepository
	logger *slog.Logger
}

func NewProductService(repo *repositories.ProductRepository, logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, logger: logger}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	s.logger.Info("Service: Getting all products")
	products, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Service: Failed to get all products", "error", err)
		return nil, err
	}
	s.logger.Info("Service: Successfully retrieved products", "count", len(products))
	return products, nil
}

func (s *ProductService) Create(data *models.Product) error {
	s.logger.Info("Service: Creating product", "name", data.Name)
	err := s.repo.Create(data)
	if err != nil {
		s.logger.Error("Service: Failed to create product", "error", err, "name", data.Name)
		return err
	}
	s.logger.Info("Service: Product created successfully", "id", data.ID)
	return nil
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	s.logger.Info("Service: Getting product by ID", "id", id)
	product, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Service: Failed to get product by ID", "error", err, "id", id)
		return nil, err
	}
	s.logger.Info("Service: Successfully retrieved product", "id", id)
	return product, nil
}

func (s *ProductService) Update(product *models.Product) error {
	s.logger.Info("Service: Updating product", "id", product.ID)
	err := s.repo.Update(product)
	if err != nil {
		s.logger.Error("Service: Failed to update product", "error", err, "id", product.ID)
		return err
	}
	s.logger.Info("Service: Product updated successfully", "id", product.ID)
	return nil
}

func (s *ProductService) Delete(id int) error {
	s.logger.Info("Service: Deleting product", "id", id)
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Service: Failed to delete product", "error", err, "id", id)
		return err
	}
	s.logger.Info("Service: Product deleted successfully", "id", id)
	return nil
}
