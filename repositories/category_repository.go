package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}


func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (Name, Description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}


func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, Name, Description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
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
			return nil, err
		}
		//bentuk mentahnya: 1, "produk A", 1000, 10

		//assign hasil slicing 
		categories = append(categories, p)
		//bentuknya {{id:1, name:"produk A", price:1000, stock:10}, {id:2, name:"produk B", price:2000, stock:20}  }
	}

	return categories, nil
}

// GetByID - ambil produk by ID
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, Name, Description FROM categories WHERE id = $1"

	var p models.Category
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4" 
	//masi HARDCODE


	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}

	return err
}
