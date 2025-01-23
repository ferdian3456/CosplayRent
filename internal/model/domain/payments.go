package domain

import "time"

type Payments struct {
	Id                        int
	Order_id                  string
	Customer_id               string
	Seller_id                 string
	Status                    string
	Method                    string
	Amount                    float64
	Payment_method            string
	Midtrans_redirect_url     string
	Midtrans_url_expired_time time.Time
	Created_at                *time.Time
	Updated_at                *time.Time
}
