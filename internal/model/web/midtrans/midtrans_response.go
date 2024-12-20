package midtrans

type MidtransResponse struct {
	Orderid     string `json:"order_id"`
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}
