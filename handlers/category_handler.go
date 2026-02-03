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

type CategoryHandler struct {
	service *services.CategoryService
	logger  *slog.Logger
}

func NewCategoryHandler(service *services.CategoryService, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{service: service, logger: logger}
	//buat object handler lalu isi dengan service yang di pass
}

// / HandleCategories - GET /api/categories
// polanya buat ini yaitu method disamping func itu receiver buat nyatain kepemilikan class yang punya
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//h itu handler w writer r request

// masuk requestnya disini nih yang parameternya write sama read
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handler: GET all categories request")
	categories, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Handler: Failed to get all categories", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
	h.logger.Info("Handler: Successfully returned all categories", "count", len(categories))
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handler: POST create category request")
	var category models.Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		h.logger.Error("Handler: Invalid request payload", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		h.logger.Error("Handler: Failed to create category", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
	h.logger.Info("Handler: Category created successfully", "id", category.ID, "name", category.Name)

}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// / getByID - GET /categories/{id}
func (h *CategoryHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Handler: Invalid category ID", "error", err, "id_str", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handler: GET category by ID request", "id", id)
	category, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error("Handler: Category not found", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
	h.logger.Info("Handler: Successfully returned category", "id", id)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Handler: Invalid category ID", "error", err, "id_str", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handler: DELETE category request", "id", id)
	//execute bussiness logic inside service module
	err = h.service.Delete(id)

	if err != nil {
		h.logger.Error("Handler: Failed to delete category (may have products)", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
	h.logger.Info("Handler: Category deleted successfully", "id", id)

}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Handler: Invalid category ID", "error", err, "id_str", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handler: PUT update category request", "id", id)
	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		h.logger.Error("Handler: Invalid request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		h.logger.Error("Handler: Failed to update category", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
	h.logger.Info("Handler: Category updated successfully", "id", id)

}
