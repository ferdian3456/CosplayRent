package review

import (
	"time"
)

type ReviewCreateRequest struct {
	User_id        string
	Costume_id     int
	Description    string
	Review_picture *string
	Rating         int
	Created_at     *time.Time
}
