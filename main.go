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
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid categories id", http.StatusBadRequest)
		return
	}

	for i, p := range categories {
		if p.ID == id {
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category deleted",
			})
			return
		}
	}
	http.Error(w, "Category belum ada", http.StatusBadRequest)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid categories id", http.StatusBadRequest)
		return
	}

	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {

			if updateCategory.Name != "" {
				categories[i].Name = updateCategory.Name
			}
			if updateCategory.Description != "" {
				categories[i].Description = updateCategory.Description
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}

	http.Error(w, "categories belum ada", http.StatusBadRequest)
}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid categories id", http.StatusBadRequest)
		return
	}

	for _, p := range categories {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "categories belum ada", http.StatusBadRequest)

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


	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var c Category
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			c.ID = len(categories) + 1
			categories = append(categories, c)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(c)
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// start server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at http://" + addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
