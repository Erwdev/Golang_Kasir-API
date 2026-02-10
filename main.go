package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/internal/logger"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"kasir-api/internal/config"
)

func main() {
	config := config.Load()

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting Kasir API", "port", config.Port)

	// Dep injection
	productRepo := repositories.NewProductRepository(db, appLogger)
	productService := services.NewProductService(productRepo, appLogger)
	productHandler := handlers.NewProductHandler(productService, appLogger)

	categoryRepo := repositories.NewCategoryRepository(db, appLogger)
	categoryService := services.NewCategoryService(categoryRepo, appLogger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, appLogger)

	transactionRepo := repositories.NewTransactionRepository(db, appLogger)
	transactionService := services.NewTransactionService(transactionRepo, appLogger)
	transactionHandler := handlers.NewTransactionHandler(transactionService, appLogger)

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "🛒 Welcome to Kasir API",
			"version":   "v0.1",
			"developer": "👨‍💻 benedictuserwdev@gmail.com",
			"endpoints": map[string]interface{}{
				"produk": map[string]string{
					"get_all":   "GET /api/produk",
					"get_by_id": "GET /api/produk/:id",
					"create":    "POST /api/produk",
					"update":    "PUT /api/produk/:id",
					"delete":    "DELETE /api/produk/:id",
				},
				"categories": map[string]string{
					"get_all":   "GET /categories",
					"get_by_id": "GET /categories/:id",
					"create":    "POST /categories",
					"update":    "PUT /categories/:id",
					"delete":    "DELETE /categories/:id",
				},
				"checkout": "POST /api/checkout",
				"metrics": "GET /metrics/db",
				"health": "GET /health",
			},
			"status": "✅ Running",
		})
	})
	

    // Add monitoring endpoint
    http.HandleFunc("/metrics/db", func(w http.ResponseWriter, r *http.Request) {
        stats := db.Stats()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "max_open_connections":   stats.MaxOpenConnections,
            "open_connections":       stats.OpenConnections,
            "in_use":                stats.InUse,
            "idle":                  stats.Idle,
            "wait_count":            stats.WaitCount,
            "wait_duration":         stats.WaitDuration.String(),
            "max_idle_closed":       stats.MaxIdleClosed,
            "max_idle_time_closed":  stats.MaxIdleTimeClosed,
            "max_lifetime_closed":   stats.MaxLifetimeClosed,
        })
    })


	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "connected"
		if err := database.HealthCheck(db); err != nil {
			dbStatus = "disconnected"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "OK",
			"database": dbStatus,
			"message":  "API running",
		})
	})

	// Produk endpoints
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// Categories endpoints
	http.HandleFunc("/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/categories/", categoryHandler.HandleCategoryByID)


	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// Start server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at http://" + addr)
	appLogger.Info("Server started successfully", "address", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		appLogger.Error("Error starting server", "error", err)
		log.Fatal("Error starting server:", err)
	}
}