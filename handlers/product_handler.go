package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
	logger  *slog.Logger
}

func NewProductHandler(service *services.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{service: service, logger: logger}
}

// / HandleProducts - GET /api/produk
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// buat single responsibility di offload ke method lain
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handler: GET all products request")
	products, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Handler: Failed to get all products", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	h.logger.Info("Handler: Successfully returned all products", "count", len(products))
}

// ...existing code...

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    var product models.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        h.logger.Error("Failed to decode request body", "error", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Debug log - lihat apa yang masuk
    h.logger.Info("Received product data", "name", product.Name, "price", product.Price, "stock", product.Stock, "category_id", product.CategoryID)

    // Validasi input
    if product.Name == "" {
        http.Error(w, "Product name is required", http.StatusBadRequest)
        return
    }
    if product.Price <= 0 {
        http.Error(w, "Product price must be greater than 0", http.StatusBadRequest)
        return
    }
    if product.Stock < 0 {
        http.Error(w, "Product stock cannot be negative", http.StatusBadRequest)
        return
    }
    if product.CategoryID <= 0 {
        http.Error(w, "Valid category_id is required", http.StatusBadRequest)
        return
    }

    err := h.service.Create(&product)
    if err != nil {
        if strings.Contains(err.Error(), "violates foreign key constraint") {
            http.Error(w, "Category ID does not exist", http.StatusBadRequest)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.logger.Info("Handler: Product created successfully", "id", product.ID, "name", product.Name)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}


//MULAI BAGIAN ENDPOINT DENGAN SLUG dengan flow selector handle

// / HandleProductByID - GET/PUT/DELETE /api/produk/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// / GetByID - GET /api/produk/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Handler: Invalid product ID", "error", err, "id_str", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handler: GET product by ID request", "id", id)
	product, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error("Handler: Product not found", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
	h.logger.Info("Handler: Successfully returned product", "id", id)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Handler: Invalid product ID", "error", err, "id_str", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handler: PUT update product request", "id", id)
	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		h.logger.Error("Handler: Invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		h.logger.Error("Handler: Failed to update product", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
	h.logger.Info("Handler: Product updated successfully", "id", id)
}

// / Delete - DELETE /api/produk/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Handler: Invalid product ID", "error", err, "id_str", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handler: DELETE product request", "id", id)
	err = h.service.Delete(id)
	if err != nil {
		h.logger.Error("Handler: Failed to delete product", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
	h.logger.Info("Handler: Product deleted successfully", "id", id)

}
