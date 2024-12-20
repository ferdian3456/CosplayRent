package user

type UserPatchRequest struct {
	Name                 *string `json:"name,omitempty" validate:"max=20"`
	Email                *string `json:"email,omitempty" validate:"max=254"`
	Address              *string `json:"address,omitempty" validate:"max=100"`
	Profile_picture      *string `json:"profile_picture,omitempty" validate:"max=255"`
	Origin_province_name *string `json:"origin_province_name,omitempty" validate:"max=30"`
	Origin_province_id   *int    `json:"origin_province_id,omitempty"`
	Origin_city_name     *string `json:"origin_city_name,omitempty" validate:"max=30"`
	Origin_city_id       *int    `json:"origin_city_id,omitempty"`
}

type TopUpEmoney struct {
	Emoney_amount float64 `validate:"required" json:"emoney_amount"`
}
