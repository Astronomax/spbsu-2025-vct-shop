package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Order struct {
	ID         int        `json:"id"`
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
	CreatedAt  string     `json:"created_at"`
}
