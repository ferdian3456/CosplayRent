package domain

import (
	"time"
)

type Order struct {
	Id              string
	User_id         string
	Seller_id       string
	Costume_id      int
	Total_amount    float64
	Status_payment  bool
	Description     string
	Status_shipping string
	Is_canceled     bool
	Created_at      *time.Time
	Updated_at      *time.Time
}
