package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
	"log/slog"
)

type CategoryRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewCategoryRepository(db *sql.DB, logger *slog.Logger) *CategoryRepository {
	return &CategoryRepository{db: db, logger: logger}
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	repo.logger.Info("Creating category", "name", category.Name, "description", category.Description)
	query := "INSERT INTO categories (Name, Description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		repo.logger.Error("Failed to create category", "error", err, "name", category.Name)
		return err
	}
	repo.logger.Info("Category created successfully", "id", category.ID, "name", category.Name)
	return nil
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	repo.logger.Info("Fetching all categories")
	query := "SELECT id, Name, Description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		repo.logger.Error("Failed to fetch categories", "error", err)
		return nil, err
	}
	defer rows.Close()
	// jalan di akhir close connect db always

	//kode di bawah buat ubah hasil query mentah jadi bentuk struct product
	categories := make([]models.Category, 0)
	for rows.Next() {
		var p models.Category

		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			repo.logger.Error("Failed to scan category", "error", err)
			return nil, err
		}
		//bentuk mentahnya: 1, "produk A", 1000, 10

		//assign hasil slicing
		categories = append(categories, p)
		//bentuknya {{id:1, name:"produk A", price:1000, stock:10}, {id:2, name:"produk B", price:2000, stock:20}  }
	}

	repo.logger.Info("Successfully fetched categories", "count", len(categories))
	return categories, nil
}

// GetByID - ambil produk by ID
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	repo.logger.Info("Fetching category by ID", "id", id)
	query := "SELECT id, Name, Description FROM categories WHERE id = $1"

	var p models.Category
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err == sql.ErrNoRows {
		repo.logger.Warn("Category not found", "id", id)
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		repo.logger.Error("Failed to fetch category by ID", "error", err, "id", id)
		return nil, err
	}

	repo.logger.Info("Successfully fetched category", "id", id, "name", p.Name)
	return &p, nil
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	repo.logger.Info("Updating category", "id", category.ID, "name", category.Name)
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"

	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		repo.logger.Error("Failed to update category", "error", err, "id", category.ID)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Failed to get rows affected", "error", err, "id", category.ID)
		return err
	}

	if rows == 0 {
		repo.logger.Warn("Category not found for update", "id", category.ID)
		return errors.New("produk tidak ditemukan")
	}

	repo.logger.Info("Category updated successfully", "id", category.ID, "name", category.Name)
	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	repo.logger.Info("Deleting category", "id", id)
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		repo.logger.Error("Failed to delete category - possibly has products referencing it", "error", err, "id", id)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Failed to get rows affected", "error", err, "id", id)
		return err
	}

	if rows == 0 {
		repo.logger.Warn("Category not found for deletion", "id", id)
		return errors.New("category tidak ditemukan")
	}

	repo.logger.Info("Category deleted successfully", "id", id)
	return nil
}
