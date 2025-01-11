package order

type OrderRequest struct {
	Seller_id             string  `validate:"required" json:"seller_id"`
	Seller_name           string  `validate:"required" json:"merchant_name"`
	Costume_id            int     `validate:"required" json:"costume_id"`
	Costume_name          string  `validate:"required" json:"costume_name"`
	Costume_category      string  `validate:"required" json:"costume_category"`
	Costume_price         float64 `validate:"required" json:"costume_price"`
	Shippment_destination string  `validate:"required" json:"shipment_destination"`
	Shipment_origin       string  `validate:"required" json:"shipment_origin"`
	TotalAmount           float64 `validate:"required" json:"total"`
	Payment_method        string  `validate:"required" json:"payment_method"`
}

type OrderEventRequest struct {
	OrderEventStatus         string `validate:"required" json:"orderevent_status"`
	OrderEventNotes          string `json:"orderevent_notes,omitempty"`
	Shipment_receipt_user_id string `json:"shipment_receipt_user_id,omitempty"`
}

type CheckBalanceWithOrderAmount struct {
	Order_amount float64 `validate:"required" json:"order_amount"`
}
