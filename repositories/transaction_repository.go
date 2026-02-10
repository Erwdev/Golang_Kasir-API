package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"log/slog"
	"sync"
)


type TransactionRepository struct {
	db *sql.DB
	logger *slog.Logger
	mu sync.Mutex
}


func NewTransactionRepository(db *sql.DB, logger *slog.Logger) *TransactionRepository {
	return &TransactionRepository{db: db, logger: logger}
}


func (repo *TransactionRepository) CreateTransaction(
	items []models.CheckoutItem,
	useLock bool,
) (*models.Transaction, error) {

	// application-level lock (optional)
	if useLock {
		repo.mu.Lock()
		defer repo.mu.Unlock()
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetails, 0, len(items))

	for _, item := range items {
		var (
			productName  string
			productPrice int
			stock        int
		)

		// lock row biar gak oversell
		err := tx.QueryRow(
			`SELECT name, price, stock
			 FROM products
			 WHERE id = $1
			 FOR UPDATE`,
			item.ProductID,
		).Scan(&productName, &productPrice, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// validasi stok
		if stock < item.Quantity {
			return nil, fmt.Errorf(
				"stok produk %s tidak cukup (tersisa %d)",
				productName,
				stock,
			)
		}

		// kurangi stok
		_, err = tx.Exec(
			`UPDATE products
			 SET stock = stock - $1
			 WHERE id = $2`,
			item.Quantity,
			item.ProductID,
		)
		if err != nil {
			return nil, err
		}

		subTotal := productPrice * item.Quantity
		totalAmount += subTotal

		details = append(details, models.TransactionDetails{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subTotal,
		})
	}

	//bakalan baca semua request semua pembelian baca db dlu baru agregasi total amountnya baru insert transaksi sama detailnya gitu, jadi gak ada masalah race condition karena kita lock row produk-nya, jadi kalau ada request lain yang mau beli produk yang sama, dia bakal nunggu sampai transaksi pertama selesai baru bisa baca stok terbaru. Jadi gak bakal oversell deh.

	// insert transaksi
	var transactionID int
	err = tx.QueryRow(
		`INSERT INTO transactions (total_amount)
		 VALUES ($1)
		 RETURNING id`,
		totalAmount,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert detail transaksi
	for i := range details {
		details[i].TransactionID = transactionID

		_, err = tx.Exec(
			`INSERT INTO transaction_details
			 (transaction_id, product_id, quantity, subtotal)
			 VALUES ($1, $2, $3, $4)`,
			details[i].TransactionID,
			details[i].ProductID,
			details[i].Quantity,
			details[i].Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
