package models

type Transaction struct {
	ID int `json:"id"`
	TotalAmount int `json:"total_amount"`
	Details []TransactionDetails `json:"details"`
}

//nested gitu ceritanya biar lebih modular mirip sama extension type gitu di typescript interface misalnya generic 

type TransactionDetails struct {
	ID int `json:"id"`
	TransactionID int `json:"transaction_id"`
	ProductID int `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity int `json:"quantity"`
	Subtotal int `json:"subtotal"`
	
}


type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
	
}
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// mirip semacem gto gitu lah kalo di nest js biar kita ada interface kalo mau ngapain gitu	
//biasain konsisten sama nama object di entity di db