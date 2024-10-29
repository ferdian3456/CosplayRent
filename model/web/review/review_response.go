package review

import (
	"time"
)

type ReviewResponse struct {
	User_id     string     `json:"user_id"`
	Costume_id  int        `json:"costume_id"`
	Description string     `json:"description"`
	Rating      int        `json:"rating"`
	Created_at  *time.Time `json:"created_at"`
}
