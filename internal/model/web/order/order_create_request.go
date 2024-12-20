package order

import (
	"time"
)

type OrderCreateRequest struct {
	Costume_id int     `validate:"required" json:"costume_id"`
	Total      float64 `validate:"required" json:"total"`
}

type DirectlyOrderToMidtrans struct {
	Id               string
	Costumer_id      string
	Costumer_email   string
	Costumer_name    string
	Seller_id        string  `validate:"required" json:"seller_id"`
	Costume_id       int     `validate:"required" json:"costume_id"`
	Costume_name     string  `validate:"required" json:"costume_name"`
	Costume_category string  `validate:"required" json:"costume_category"`
	Costume_price    float64 `validate:"required" json:"costume_price"`
	Merchant_name    string  `validate:"required" json:"merchant_name"`
	TotalAmount      float64 `validate:"required" json:"total"`
	Created_at       *time.Time
	Updated_at       *time.Time
}

type CheckBalanceWithOrderAmount struct {
	Order_amount float64 `validate:"required" json:"order_amount"`
}
