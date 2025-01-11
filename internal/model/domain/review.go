package domain

import (
	"time"
)

type Review struct {
	Id             int
	User_id        string
	Order_id       string
	Costume_id     int
	Description    string
	Review_picture string
	Rating         int
	Created_at     *time.Time
	Updated_at     *time.Time
}
