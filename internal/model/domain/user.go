package domain

import (
	"time"
)

type User struct {
	Id                   string
	Name                 string
	Email                string
	Address              string
	Password             string
	Profile_picture      string
	Origin_province_name string
	Origin_province_id   int
	Origin_city_name     string
	Origin_city_id       int
	Created_at           *time.Time
	Updated_at           *time.Time
}
