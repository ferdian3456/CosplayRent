package review

type OwnReviewByReviewID struct {
	User_id        string
	Costume_id     int
	Username       string  `json:"name"`
	Costume_name   string  `json:"costume_name"`
	Rating         int     `json:"rating"`
	Description    string  `json:"description"`
	Review_picture *string `json:"-"`
	Created_at     string  `json:"created_at"`
	Updated_at     string  `json:"updated_at"`
}
