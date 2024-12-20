package user

type UserPatchRequest struct {
	Name                 *string `json:"name" validate:"max=20"`
	Email                *string `json:"email" validate:"max=254"`
	Address              *string `json:"address" validate:"max=100"`
	Profile_picture      *string `json:"profile_picture"`
	Origin_province_name *string `json:"origin_province_name" validate:"max=30"`
	Origin_province_id   *int    `json:"origin_province_id"`
	Origin_city_name     *string `json:"origin_city_name" validate:"max=30"`
	Origin_city_id       *int    `json:"origin_city_id"`
}

type TopUpEmoney struct {
	Emoney_amount float64 `validate:"required" json:"emoney_amount"`
}
