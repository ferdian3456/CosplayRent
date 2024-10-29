package domain

import (
	"time"
)

type Review struct {
	User_id     string
	Costume_id  int
	Description string
	Rating      int
	Created_at  *time.Time
}
