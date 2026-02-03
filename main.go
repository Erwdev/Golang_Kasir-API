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

// Config hanya pegang PORT dan DB_CONN.
// Railway set env di OS level, bukan di file .env,
// jadi kita baca os.Getenv() langsung sebagai prioritas utama.
type Config struct {
	Port   string
	DBConn string
}

// loadConfig baca config dengan urutan prioritas:
//  1. OS environment variable (Railway set di sini)
//  2. .env file (untuk local development)
//  3. Default value
func loadConfig() Config {
	// Coba baca .env untuk local dev — kalau file-nya ga ada, ga error
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		viper.ReadInConfig() // error di-ignore, wajar kalau ga ada
	}

	// Baca dari viper sebagai fallback (dari .env)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := Config{}

	// PORT: OS env > .env > default "8080"
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	} else if port := viper.GetString("PORT"); port != "" {
		cfg.Port = port
	} else {
		cfg.Port = "8080"
	}

	// DB_CONN: OS env > .env
	// Railway biasanya kasih DATABASE_URL langsung,
	// tapi kalau lo pake env name DB_CONN di Railway dashboard, itu juga fine.
	if dbConn := os.Getenv("DB_CONN"); dbConn != "" {
		cfg.DBConn = dbConn
	} else if dbConn := os.Getenv("DATABASE_URL"); dbConn != "" {
		// fallback ke DATABASE_URL yang Railway auto-generate
		cfg.DBConn = dbConn
	} else {
		cfg.DBConn = viper.GetString("DB_CONN")
	}

	return cfg
}

func main() {
	config := loadConfig()

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
				"health": "GET /health",
			},
			"status": "✅ Running",
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

	// Start server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at http://" + addr)
	appLogger.Info("Server started successfully", "address", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		appLogger.Error("Error starting server", "error", err)
		log.Fatal("Error starting server:", err)
	}
}