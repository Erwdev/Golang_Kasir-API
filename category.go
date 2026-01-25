// buat misahin data aja biar rapi n bisa generate banyak sama gpt isinya

package main

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Beras Premium 5kg", Description: "Beras putih kualitas premium kemasan 5 kilogram"},
	{ID: 2, Name: "Minyak Goreng 2L", Description: "Minyak goreng sawit kemasan botol 2 liter"},
	{ID: 3, Name: "Gula Pasir 1kg", Description: "Gula pasir putih kemasan 1 kilogram"},
	{ID: 4, Name: "Susu UHT Coklat", Description: "Susu UHT rasa coklat kemasan 1 liter"},
	{ID: 5, Name: "Mie Instan Goreng", Description: "Mie instan goreng rasa original"},
	{ID: 6, Name: "Kopi Bubuk 200g", Description: "Kopi bubuk robusta kemasan 200 gram"},
	{ID: 7, Name: "Teh Celup 25pcs", Description: "Teh celup hitam isi 25 kantong"},
	{ID: 8, Name: "Sabun Mandi Cair", Description: "Sabun mandi cair dengan aroma segar"},
	{ID: 9, Name: "Pasta Gigi 120g", Description: "Pasta gigi dengan perlindungan gigi berlubang"},
	{ID: 10, Name: "Air Mineral 600ml", Description: "Air mineral dalam kemasan botol 600 ml"},
}


// jangan lupa pake package main biar bisa ke detect satu namespace dan di link pas compiling dengan main.go sebagai entry point 
//jangan lupa compile semuanya pake go run . (all) aja dlu 