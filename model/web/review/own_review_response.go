package review

type OwnReviewResponse struct {
	Id              int     `json:"id"`
	Userid          string  `json:"-"`
	Costume_Id      int     `json:"-"`
	Review_picture  string  `json:"-"`
	Costume_name    string  `json:"costume_name"`
	Costume_picture *string `json:"costume_picture"`
	Ukuran          *string `json:"ukuran"`
	Description     string  `json:"description"`
	Rating          int     `json:"rating"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}
