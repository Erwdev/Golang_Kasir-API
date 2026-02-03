package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
	"log/slog"
)

type ProductRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewProductRepository(db *sql.DB, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

func (repo *ProductRepository) Create(product *models.Product) error {


	repo.logger.Info("Creating product", "name", product.Name, "price", product.Price, "stock", product.Stock, "category_id", product.CategoryID)
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		repo.logger.Error("Failed to create product", "error", err, "name", product.Name)
		return err
	}
	repo.logger.Info("Product created successfully", "id", product.ID, "name", product.Name)
	return nil
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	repo.logger.Info("Fetching all products")
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name
		FROM products p 
		LEFT JOIN categories c ON p.category_id = c.id
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		repo.logger.Error("Failed to fetch products", "error", err)
		return nil, err
	}
	defer rows.Close()
	// jalan di akhir close connect db always

	//kode di bawah buat ubah hasil query mentah jadi bentuk struct product
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
		if err != nil {
			repo.logger.Error("Failed to scan product", "error", err)
			return nil, err
		}
		//bentuk mentahnya: 1, "produk A", 1000, 10

		//assign hasil slicing
		products = append(products, p)
		//bentuknya {{id:1, name:"produk A", price:1000, stock:10}, {id:2, name:"produk B", price:2000, stock:20}  }
	}

	repo.logger.Info("Successfully fetched products", "count", len(products))
	return products, nil
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	repo.logger.Info("Fetching product by ID", "id", id)
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name 
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
	if err == sql.ErrNoRows {
		repo.logger.Warn("Product not found", "id", id)
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		repo.logger.Error("Failed to fetch product by ID", "error", err, "id", id)
		return nil, err
	}

	repo.logger.Info("Successfully fetched product", "id", id, "name", p.Name)
	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	repo.logger.Info("Updating product", "id", product.ID, "name", product.Name)
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	//masi HARDCODE

	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		repo.logger.Error("Failed to update product", "error", err, "id", product.ID)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Failed to get rows affected", "error", err, "id", product.ID)
		return err
	}

	if rows == 0 {
		repo.logger.Warn("Product not found for update", "id", product.ID)
		return errors.New("produk tidak ditemukan")
	}

	repo.logger.Info("Product updated successfully", "id", product.ID, "name", product.Name)
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	repo.logger.Info("Deleting product", "id", id)
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		repo.logger.Error("Failed to delete product", "error", err, "id", id)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Failed to get rows affected", "error", err, "id", id)
		return err
	}

	if rows == 0 {
		repo.logger.Warn("Product not found for deletion", "id", id)
		return errors.New("produk tidak ditemukan")
	}

	repo.logger.Info("Product deleted successfully", "id", id)
	return err
}
