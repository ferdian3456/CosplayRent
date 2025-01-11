package order

import (
	"github.com/google/uuid"
)

type OrderResponse struct {
	Id              uuid.UUID `json:"id"`
	Customer_id     uuid.UUID `json:"customer_id"`
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

type OrderForReviewsResponse struct {
	Id         string `json:"id"`
	Seller_id  string `json:"seller_id"`
	Custome_Id int    `json:"costume_id"`
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
	Shipment_destination     string  `json:"shipment_destination"`
	Costumer_identity_card   string  `json:"costumer_identity_card"`
	Shipment_receipt_user_id string  `json:"shipment_receipt_user_id,omitempty"`
	Shipment_notes           string  `json:"shipment_notes,omitempty"`
}

type GetUserOrderDetailResponse struct {
	Costume_picture          *string `json:"costume_picture"`
	Costume_name             string  `json:"costume_name"`
	Costume_price            float64 `json:"costume_price"`
	Costume_size             *string `json:"costume_size"`
	Seller_name              string  `json:"seller_name"`
	Seller_address           *string `json:"seller_address"`
	Seller_response          string  `json:"seller_response"`
	Shipment_receipt_user_id string  `json:"shipment_receipt_user_id,omitempty"`
	Shipment_notes           string  `json:"shipment_notes,omitempty"`
}

type CheckBalanceWithOrderAmountReponse struct {
	Status_to_order string `json:"status_to_order"`
}

type PaymentTransationForOrderResponse struct {
	Order_id                           string  `json:"-"`
	Costume_id                         int     `json:"-"`
	Payment_status                     string  `json:"payment_status"`
	Costume_name                       string  `json:"costume_name"`
	Costume_picture                    *string `json:"costume_picture"`
	Costume_price                      float64 `json:"costume_price"`
	Costume_size                       string  `json:"costume_size"`
	Midtrans_redirect_url_expired_time string  `json:"midtrans_expired_time"`
	Payment_id                         int     `json:"payment_id"`
}

type PaymentInfo struct {
	Payment_amount                     float64 `json:"payment_amount"`
	Status                             string  `json:"payment_status"`
	Midtrans_redirect_url              string  `json:"midtrans_redirect_url"`
	Midtrans_redirect_url_expired_time string  `json:"midtrans_expired_time"`
}
