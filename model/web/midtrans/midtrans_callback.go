package midtrans

type MidtransCallback struct {
	Status_Code       string `json:"status_code"`
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionID     string `json:"transaction_id"`
	SignatureKey      string `json:"signature_key"`
}
