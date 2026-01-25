package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/") //trim prefix yang ditentuin hard coded

	id, err := strconv.Atoi(idStr) //convert dari string ke int
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return

	}
	//cari brute force
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product belum ada", http.StatusBadRequest)

}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id //biar tetep sama
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	http.Error(w, "Product belum ada", http.StatusBadRequest)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...) //mirip spread operator di js??
			//basically masukin lagi by slicing bagian kue dari sisi kiri sama kanan index yang jadi target wkkwwk mirip python

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Product deleted",
			}) //delete cukup konfirmasi aja
			return
		}
	}
	http.Error(w, "Product belum ada", http.StatusBadRequest)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "🛒 Welcome to Kasir API",
			"version":   "v0.1",
			"developer": "👨‍💻 benedictuserwdev@gmail.com",
			"endpoints": map[string]string{
				"health":    "GET /health",
				"get_all":   "GET /api/produk",
				"get_by_id": "GET /api/produk/:id",
				"create":    "POST /api/produk",
				"update":    "PUT /api/produk/:id",
				"delete":    "DELETE /api/produk/:id",
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

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	// bawah handler buat tanpa slash endpoint karena lebih simple katanya
	// yang atas buat handling dengan slug parameter id

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)

		} else if r.Method == "POST" {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid Request Body", http.StatusBadRequest)
				return
			}
			// pake pointer buat arahin decoder simpen di alamat situ

			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) //201 status code
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	} //jangan lupa masukin return nil di main biar ga error dan handling
}
