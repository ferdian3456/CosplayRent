package wishlist

type WishListResponses struct {
	CostumeId       int     `json:"costume_id"`
	Costume_name    string  `json:"costume_name"`
	Costume_picture *string `json:"costume_picture"`
	Costume_price   float64 `json:"costume_price"`
	Costume_size    string  `json:"costume_size"`
}
