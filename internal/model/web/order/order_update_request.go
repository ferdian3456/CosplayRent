package order

import "time"

type OrderUpdateRequest struct {
	StatusOrder string `validate:"required" json:"status_order"`
	Description string `json:"description"`
	Updated_at  *time.Time
}
