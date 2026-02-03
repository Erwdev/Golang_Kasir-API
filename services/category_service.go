package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"log/slog"
)

type CategoryService struct {
	repo   *repositories.CategoryRepository
	logger *slog.Logger
}

func NewCategoryService(repo *repositories.CategoryRepository, logger *slog.Logger) *CategoryService {
	return &CategoryService{repo: repo, logger: logger}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	s.logger.Info("Service: Getting all categories")
	categories, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Service: Failed to get all categories", "error", err)
		return nil, err
	}
	s.logger.Info("Service: Successfully retrieved categories", "count", len(categories))
	return categories, nil
}

func (s *CategoryService) Create(data *models.Category) error {
	s.logger.Info("Service: Creating category", "name", data.Name)
	err := s.repo.Create(data)
	if err != nil {
		s.logger.Error("Service: Failed to create category", "error", err, "name", data.Name)
		return err
	}
	s.logger.Info("Service: Category created successfully", "id", data.ID)
	return nil
}

func (s *CategoryService) Delete(id int) error {
	s.logger.Info("Service: Deleting category", "id", id)
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Service: Failed to delete category (may have foreign key constraint)", "error", err, "id", id)
		return err
	}
	s.logger.Info("Service: Category deleted successfully", "id", id)
	return nil
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	s.logger.Info("Service: Getting category by ID", "id", id)
	category, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Service: Failed to get category by ID", "error", err, "id", id)
		return nil, err
	}
	s.logger.Info("Service: Successfully retrieved category", "id", id)
	return category, nil
}

func (s *CategoryService) Update(data *models.Category) error {
	s.logger.Info("Service: Updating category", "id", data.ID)
	err := s.repo.Update(data)
	if err != nil {
		s.logger.Error("Service: Failed to update category", "error", err, "id", data.ID)
		return err
	}
	s.logger.Info("Service: Category updated successfully", "id", data.ID)
	return nil
}
