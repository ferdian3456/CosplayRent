package order

import (
	"time"
)

type OrderCreateRequest struct {
	Id         string
	User_id    string  `validate:"required" json:"user_id"`
	Costume_id int     `validate:"required" json:"costume_id"`
	Total      float64 `validate:"required" json:"total"`
	Created_at time.Time
}
