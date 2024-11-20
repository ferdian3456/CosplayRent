package review

import "time"

type OwnReviewResponse struct {
	Id              int `json:"id"`
	Userid          string
	Costume_Id      int
	Review_picture  string
	Costume_name    string     `json:"costume_name"`
	Costume_picture *string    `json:"costume_picture"`
	Ukuran          *string    `json:"ukuran"`
	Description     string     `json:"description"`
	Rating          int        `json:"rating"`
	Created_at      *time.Time `json:"created_at"`
}
