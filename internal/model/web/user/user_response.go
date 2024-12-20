package user

type UserResponse struct {
	Id                   string  `json:"id"`
	Name                 string  `json:"name"`
	Email                string  `json:"email"`
	Address              *string `json:"address"`
	Profile_picture      *string `json:"profile_picture"`
	Origin_province_name *string `json:"origin_province_name"`
	Origin_province_id   *int    `json:"origin_province_id"`
	Origin_city_name     *string `json:"origin_city_name"`
	Origin_city_id       *int    `json:"origin_city_id"`
	Created_at           string  `json:"created_at"`
	Updated_at           string  `json:"updated_at"`
}

type IdentityCardResponse struct {
	IdentityCard_picture string `json:"identitycard_picture"`
}

type UserEmoneyResponse struct {
	Emoney_amont      float64 `json:"emoney_amount"`
	Emoney_updated_at string  `json:"emoney_updated_at"`
}

type UserEMoneyTransactionHistory struct {
	Transaction_amount float64 `json:"transaction_amount"`
	Transaction_type   string  `json:"transaction_type"`
	Transaction_date   string  `json:"transaction_date"`
}

type CheckUserStatusResponse struct {
	User_id string `json:"user_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type SellerAddressResponse struct {
	Seller_name                 string  `json:"seller_name"`
	Seller_origin_province_name *string `json:"seller_origin_province_name"`
	Seller_origin_province_id   *int    `json:"seller_origin_province_id"`
	Seller_origin_city_name     *string `json:"seller_origin_city_name"`
	Seller_origin_city_id       *int    `json:"seller_origin_city_id"`
}
