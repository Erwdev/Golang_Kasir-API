// buat misahin data aja biar rapi n bisa generate banyak sama gpt isinya

package main

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
    {ID: 1, Name: "Makanan", Description: "Kategori untuk makanan ringan dan berat"},
    {ID: 2, Name: "Minuman", Description: "Kategori untuk minuman kemasan"},
    {ID: 3, Name: "Bumbu Dapur", Description: "Kategori untuk bumbu dan penyedap"},
    {ID: 4, Name: "Sembako", Description: "Kategori untuk kebutuhan pokok"},
    {ID: 5, Name: "Snack", Description: "Kategori untuk camilan dan makanan ringan"},
    {ID: 6, Name: "Toiletries", Description: "Kategori untuk keperluan mandi dan kebersihan"},
    {ID: 7, Name: "Peralatan Dapur", Description: "Kategori untuk alat-alat dapur"},
    {ID: 8, Name: "Frozen Food", Description: "Kategori untuk makanan beku"},
    {ID: 9, Name: "Produk Kesehatan", Description: "Kategori untuk vitamin dan obat-obatan"},
    {ID: 10, Name: "Elektronik", Description: "Kategori untuk barang elektronik"},
}

// jangan lupa pake package main biar bisa ke detect satu namespace dan di link pas compiling dengan main.go sebagai entry point 
//jangan lupa compile semuanya pake go run . (all) aja dlu 