package domain

import "time"

type Midtrans struct {
	Order_id       string
	Order_amount   float64
	TopUpUser_id   string
	OrderBuyer_id  string
	OrderSeller_id string
	Updated_at     *time.Time
}
