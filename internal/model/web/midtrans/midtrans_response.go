package midtrans

type MidtransResponse struct {
	Orderid            string `json:"order_id"`
	Token              string `json:"token"`
	RedirectUrl        string `json:"redirect_url"`
	MidtransExpired    string `json:"midtrans_expired_time"`
	MidtransCreated_at string `json:"midtrans_created_at"`
}
