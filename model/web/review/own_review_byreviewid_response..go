package review

import "time"

type OwnReviewByReviewID struct {
	User_id      string
	Costume_id   int
	Username     string     `json:"name"`
	Costume_name string     `json:"costume_name"`
	Rating       int        `json:"rating"`
	Description  string     `json:"description"`
	Created_at   *time.Time `json:"created_at"`
}
