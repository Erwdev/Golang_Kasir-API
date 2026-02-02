package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
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
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	//dep injection 
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	

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
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API running",
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

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
