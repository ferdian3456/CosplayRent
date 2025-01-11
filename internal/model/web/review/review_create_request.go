package review

type ReviewCreateRequest struct {
	Order_id       string  `validate:"required" json:"order_id"`
	Customer_id    string  `json:"-"`
	Costume_id     int     `validate:"required" json:"costume_id"`
	Description    string  `validate:"required" json:"description"`
	Review_picture *string `validate:"required" json:"review_picture"`
	Rating         int     `validate:"required" json:"rating"`
}
