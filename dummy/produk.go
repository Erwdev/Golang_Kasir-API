// buat misahin data aja biar rapi n bisa generate banyak sama gpt isinya

package main

type Produk struct {
	ID int `json:"id"`
	Nama string `json:"nama"`
	Harga int `json:"harga"`
	Stok int `json:"stok"`
}

var produk = []Produk{
	
	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 50},
	{ID: 2, Nama: "Indomie Godog", Harga: 3500, Stok: 40},
	{ID: 3, Nama: "Indomie Rendang", Harga: 4000, Stok: 30},
	{ID: 4, Nama: "Indomie Soto", Harga: 3500, Stok: 25},
	{ID: 5, Nama: "Indomie Ayam Bawang", Harga: 3500, Stok: 60},

	{ID: 6, Nama: "Aqua 600ml", Harga: 3000, Stok: 100},
	{ID: 7, Nama: "Aqua 1500ml", Harga: 6000, Stok: 70},
	{ID: 8, Nama: "Vit 600ml", Harga: 2500, Stok: 90},
	{ID: 9, Nama: "Vit 1000ml", Harga: 3000, Stok: 80},
	{ID: 10, Nama: "Le Minerale 600ml", Harga: 3000, Stok: 85},

	{ID: 11, Nama: "Teh Pucuk", Harga: 4000, Stok: 45},
	{ID: 12, Nama: "Teh Botol Sosro", Harga: 5000, Stok: 55},
	{ID: 13, Nama: "Nu Green Tea", Harga: 4500, Stok: 35},
	{ID: 14, Nama: "Frestea", Harga: 4000, Stok: 40},
	{ID: 15, Nama: "Teh Kotak", Harga: 4500, Stok: 50},

	{ID: 16, Nama: "Kopi Kapal Api", Harga: 2000, Stok: 120},
	{ID: 17, Nama: "Good Day Cappuccino", Harga: 2500, Stok: 90},
	{ID: 18, Nama: "Torabika Susu", Harga: 2000, Stok: 110},
	{ID: 19, Nama: "ABC Kopi Susu", Harga: 2000, Stok: 100},
	{ID: 20, Nama: "Nescafe Classic", Harga: 5000, Stok: 60},

	{ID: 21, Nama: "Kecap Bango", Harga: 12000, Stok: 40},
	{ID: 22, Nama: "Kecap ABC", Harga: 9000, Stok: 50},
	{ID: 23, Nama: "Saos Sambal ABC", Harga: 8000, Stok: 45},
	{ID: 24, Nama: "Saos Tomat ABC", Harga: 8000, Stok: 35},
	{ID: 25, Nama: "Sambal Jawara", Harga: 10000, Stok: 30},

	{ID: 26, Nama: "Gula Pasir 1kg", Harga: 14000, Stok: 70},
	{ID: 27, Nama: "Gula Merah", Harga: 15000, Stok: 40},
	{ID: 28, Nama: "Tepung Terigu 1kg", Harga: 12000, Stok: 60},
	{ID: 29, Nama: "Tepung Beras", Harga: 11000, Stok: 50},
	{ID: 30, Nama: "Minyak Goreng 1L", Harga: 18000, Stok: 80},

	{ID: 31, Nama: "Minyak Goreng 2L", Harga: 35000, Stok: 45},
	{ID: 32, Nama: "Mentega Blue Band", Harga: 9000, Stok: 55},
	{ID: 33, Nama: "Margarin Palmia", Harga: 8500, Stok: 60},
	{ID: 34, Nama: "Susu Ultra Milk", Harga: 6000, Stok: 75},
	{ID: 35, Nama: "Susu Indomilk", Harga: 5500, Stok: 70},

	{ID: 36, Nama: "Roti Tawar Sari Roti", Harga: 15000, Stok: 30},
	{ID: 37, Nama: "Roti Sobek", Harga: 12000, Stok: 25},
	{ID: 38, Nama: "Biskuit Roma", Harga: 8000, Stok: 65},
	{ID: 39, Nama: "Biskuit Marie", Harga: 7000, Stok: 60},
	{ID: 40, Nama: "Chocolatos", Harga: 2000, Stok: 150},

	{ID: 41, Nama: "SilverQueen", Harga: 12000, Stok: 40},
	{ID: 42, Nama: "Delfi Chocolate", Harga: 10000, Stok: 35},
	{ID: 43, Nama: "Qtela Original", Harga: 7000, Stok: 50},
	{ID: 44, Nama: "Chitato", Harga: 8000, Stok: 45},
	{ID: 45, Nama: "Taro Net", Harga: 7500, Stok: 55},

	{ID: 46, Nama: "Mi Telur", Harga: 6000, Stok: 60},
	{ID: 47, Nama: "Beras Ramos 5kg", Harga: 65000, Stok: 20},
	{ID: 48, Nama: "Beras Pandan Wangi", Harga: 70000, Stok: 15},
	{ID: 49, Nama: "Garam Dapur", Harga: 4000, Stok: 90},
	{ID: 50, Nama: "Kaldu Ayam", Harga: 5000, Stok: 85},
}


// jangan lupa pake package main biar bisa ke detect satu namespace dan di link pas compiling dengan main.go sebagai entry point 
//jangan lupa compile semuanya pake go run . (all) aja dlu 