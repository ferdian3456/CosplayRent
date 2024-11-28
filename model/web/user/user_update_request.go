package user

import "time"

type UserUpdateRequest struct {
	Id                   string
	Name                 *string `json:"name"`
	Email                *string `json:"email"`
	Address              *string `json:"address"`
	Profile_picture      *string `json:"profile_picture"`
	Origin_province_name *string `json:"origin_province_name"`
	Origin_city_name     *string `json:"origin_city_name"`
	Created_at           *time.Time
	Update_at            *time.Time
}

type TopUpEmoney struct {
	Emoney_amont float64 `json:"emoney_amount"`
}
