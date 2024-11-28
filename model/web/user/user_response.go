package user

type UserResponse struct {
	Id                   string  `json:"id"`
	Name                 string  `json:"name"`
	Email                string  `json:"email"`
	Address              *string `json:"address"`
	Profile_picture      *string `json:"profile_picture"`
	Origin_city_name     *string `json:"origin_city_name"`
	Origin_province_name *string `json:"origin_province_name"`
	Created_at           string  `json:"created_at"`
	Updated_at           string  `json:"updated_at"`
}
