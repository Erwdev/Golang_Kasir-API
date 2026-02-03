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
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
	}

	// Setelah ReadInConfig, override balik dari OS env
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("PORT", port) // Set punya prioritas tertinggi
	}
	if dbConn := os.Getenv("DB_CONN"); dbConn != "" {
		viper.Set("DB_CONN", dbConn)
	}

	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	config := Config{
		Port:   port,
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting Kasir API", "port", config.Port)

	//dep injection
	productRepo := repositories.NewProductRepository(db, appLogger)
	productService := services.NewProductService(productRepo, appLogger)
	productHandler := handlers.NewProductHandler(productService, appLogger)

	categoryRepo := repositories.NewCategoryRepository(db, appLogger)
	categoryService := services.NewCategoryService(categoryRepo, appLogger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, appLogger)

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

				"health": "GET /health",
			},
			"status": "✅ Running",
		})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// json.NewEncoder(w).Encode(map[string]string{
		// 	"status":  "OK",
		// 	"message": "API running",
		// })
		dbStatus := "connected"
		if err := database.HealthCheck(db); err != nil {
			dbStatus = "disconnected"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "OK",
			"database": dbStatus,
			"message":  "API running",
		})
	})

	//buat endpoint tanpa slug
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	//endpoint dengan slug
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	//buat endpoint tanpa slug
	http.HandleFunc("/categories", categoryHandler.HandleCategories)
	//endpoint dengan slug
	http.HandleFunc("/categories/", categoryHandler.HandleCategoryByID)

	// start server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at http://" + addr)
	appLogger.Info("Server started successfully", "address", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		appLogger.Error("Error starting server", "error", err)
		fmt.Println("Error starting server:", err)
	}
}
