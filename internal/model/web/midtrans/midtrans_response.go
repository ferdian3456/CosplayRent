package midtrans

import "github.com/google/uuid"

type MidtransResponse struct {
	Orderid     uuid.UUID `json:"order_id"`
	Token       string    `json:"token"`
	RedirectUrl string    `json:"redirect_url"`
}
