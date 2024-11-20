package review

import (
	"time"
)

type ReviewUpdateRequest struct {
	ReviewId       int
	Review_picture *string `json:"review_picture"`
	Description    *string `json:"description"`
	Rating         *string `json:"rating"`
	Updated_at     *time.Time
}
