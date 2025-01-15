package review

type ReviewResponse struct {
	Costume_id      string  `json:"costume_id"`
	User_id         string  `json:"user_id"`
	Name            string  `json:"name"`
	Profile_picture *string `json:"profile_picture"`
	Review_picture  *string `json:"review_picture"`
	Description     string  `json:"description"`
	Rating          int     `json:"rating"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}

type UserReviewResponse struct {
	Id              string  `json:"-"`
	Review_picture  *string `json:"review_picture"`
	Seller_id       string  `json:"-"`
	Custome_Id      int     `json:"-"`
	Order_id        string  `json:"order_id"`
	Seller_name     string  `json:"seller_name"`
	Costume_name    string  `json:"costume_name"`
	Costume_picture string  `json:"costume_picture"`
	Costume_size    string  `json:"costume_size"`
	Costume_weight  int     `json:"costume_weight"`
}

type UserReviewDetailByIdResponse struct {
	Seller_name          string  `json:"seller_name"`
	Costume_name         string  `json:"costume_name"`
	Costume_picture      string  `json:"costume_picture"`
	Costume_size         string  `json:"costume_size"`
	Costume_weight       int     `json:"costume_weight"`
	Costume_material     string  `json:"costume_material"`
	Order_amount         float64 `json:"order_amount"`
	Shipment_destination string  `json:"shipment_destination"`
	Shipment_origin      string  `json:"shipment_origin"`
	Picture              *string `json:"review_picture"`
	Rating               *int    `json:"rating"`
	Description          *string `json:"description"`
	Created_at           *string `json:"created_at"`
	Updated_at           *string `json:"updated_at"`
}

type ListOfNonReviewedOrder struct {
	Order_id        string  `json:"order_id"`
	Order_date      string  `json:"order_date"`
	Costume_id      int     `json:"costume_id"`
	Costume_name    string  `json:"costume_name"`
	Costume_picture *string `json:"costume_picture"`
	Costume_size    string  `json:"costume_size"`
	Costume_price   float64 `json:"costume_price"`
}

type ListOfReviewedOrder struct {
	Order_id           string  `json:"order_id"`
	Review_rating      int     `json:"review_rating"`
	Review_description string  `json:"review_description"`
	Review_date        string  `json:"review_date"`
	Costume_id         int     `json:"costume_id"`
	Costume_name       string  `json:"costume_name"`
	Costume_picture    *string `json:"costume_picture"`
}

type ReviewInfo struct {
	Review_id          int    `json:"review_id"`
	Review_rating      int    `json:"review_rating"`
	Review_picture     string `json:"review_picture"`
	Review_description string `json:"review_description"`
}
