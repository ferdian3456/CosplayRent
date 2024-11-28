package order

import (
	"github.com/google/uuid"
)

type OrderResponse struct {
	Id              uuid.UUID `json:"id"`
	User_id         uuid.UUID `json:"user_id"`
	Seller_id       uuid.UUID `json:"seller_id"`
	Costume_id      int       `json:"costume_id"`
	Shipping_id     int       `json:"shipping_id"`
	Total           float64   `json:"total"`
	Status_payment  bool      `json:"status_payment"`
	Status_shipping bool      `json:"status_shipping"`
	Is_canceled     bool      `json:"is_canceled"`
	Created_at      string    `json:"created_at"`
	Updated_at      string    `json:"updated_at"`
}
