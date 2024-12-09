package order

import (
	"github.com/google/uuid"
)

type OrderResponse struct {
	Id              uuid.UUID `json:"id"`
	User_id         uuid.UUID `json:"user_id"`
	Seller_id       uuid.UUID `json:"seller_id"`
	Costume_id      int       `json:"costume_id"`
	Description     *string   `json:"description"`
	Shipping_id     int       `json:"shipping_id"`
	Total           float64   `json:"total"`
	Status_payment  bool      `json:"status_payment"`
	Status_shipping bool      `json:"status_shipping"`
	Is_canceled     bool      `json:"is_canceled"`
	Created_at      string    `json:"created_at"`
	Updated_at      string    `json:"updated_at"`
}

type AllSellerOrderResponse struct {
	Id              uuid.UUID `json:"id"`
	Status          string    `json:"status_order"`
	Costume_id      int
	Costume_name    string  `json:"costume_name"`
	Costume_price   float64 `json:"costume_price"`
	Costume_size    *string `json:"costume_size"`
	Costume_picture *string `json:"costume_picture"`
	Total           float64 `json:"total"`
	Updated_at      string  `json:"updated_at"`
}

type AllUserOrderResponse struct {
	Id              uuid.UUID `json:"id"`
	Status          string    `json:"status_order"`
	Costume_id      int
	Costume_name    string  `json:"costume_name"`
	Costume_price   float64 `json:"costume_price"`
	Costume_size    *string `json:"costume_size"`
	Costume_picture *string `json:"costume_picture"`
	Total           float64 `json:"total"`
	Updated_at      string  `json:"updated_at"`
}

type OrderDetailByOrderIdResponse struct {
	Costume_picture          *string `json:"costume_picture"`
	Costume_name             string  `json:"costume_name"`
	Costume_price            float64 `json:"costume_price"`
	Costume_size             *string `json:"costume_size"`
	Costumer_name            string  `json:"costumer_name"`
	Costumer_address         *string `json:"costumer_address"`
	Costumer_origin_province *string `json:"costumer_origin_province"`
	Costumer_origin_city     *string `json:"costumer_origin_city"`
	Costumer_identity_card   string  `json:"costumer_identity_card"`
}

type GetUserOrderDetailResponse struct {
	Costume_picture        *string `json:"costume_picture"`
	Costume_name           string  `json:"costume_name"`
	Costume_price          float64 `json:"costume_price"`
	Costume_size           *string `json:"costume_size"`
	Seller_name            string  `json:"seller_name"`
	Seller_address         *string `json:"seller_address"`
	Seller_origin_province *string `json:"seller_origin_province"`
	Seller_origin_city     *string `json:"seller_origin_city"`
	Seller_response        *string `json:"seller_response"`
}

type CheckBalanceWithOrderAmountReponse struct {
	Status_to_order string `json:"status_to_order"`
}
