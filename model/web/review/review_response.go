package review

import (
	"time"
)

type ReviewResponse struct {
	User_id         string     `json:"-"`
	Costume_id      int        `json:"-"`
	Name            string     `json:"name"`
	Profile_picture *string    `json:"profile_picture"`
	Description     string     `json:"description"`
	Rating          int        `json:"rating"`
	Created_at      *time.Time `json:"created_at"`
}
