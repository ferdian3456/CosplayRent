package order

import (
	"github.com/google/uuid"
	"time"
)

type OrderResponse struct {
	Id             uuid.UUID `json:"id"`
	User_id        uuid.UUID `json:"user_id"`
	Costume_id     int       `json:"costume_id"`
	Shipping_id    int       `json:"shipping_id"`
	Total          float64   `json:"total"`
	Status_payment bool      `json:"status_payment"`
	Is_canceled    bool      `json:"is_canceled"`
	Created_at     time.Time `json:"created_at"`
}
