package review

import (
	"time"
)

type ReviewCreateRequest struct {
	User_id     string `validate:"required,min=36,max=36" json:"user_id"`
	Costume_id  int    `validate:"required,max=36" json:"costume_id"`
	Description string `validate:"required,min=5,max=1000" json:"description"`
	Rating      int    `validate:"required" json:"rating"`
	Created_at  *time.Time
}
