package domain

import (
  "time"
)

type Order struct {
	Id             string
	User_id        string
	Costume_id     int
	Shipping_id    int
	Total          float64
	Status_payment bool
	Is_canceled    bool
	Created_at     *time.Time
}
