package domain

import (
	"time"
)

type Order struct {
	Id                   string
	Seller_id            string
	Costumer_id          string
	Costume_id           int
	Total_amount         float64
	Shipment_origin      string
	Shipment_destination string
	Created_at           *time.Time
	Updated_at           *time.Time
}

type OrderToMidtrans struct {
	Id               string
	Seller_id        string
	Seller_name      string
	Costumer_id      string
	Costumer_name    string
	Costumer_email   string
	Costume_id       int
	Costume_name     string
	Costume_category string
	Costume_price    float64
	Total_amount     float64
	Created_at       time.Time
}
