package review

type ReviewResponse struct {
	User_id         string  `json:"-"`
	Costume_id      int     `json:"-"`
	Name            string  `json:"name"`
	Profile_picture *string `json:"profile_picture"`
	Description     string  `json:"description"`
	Rating          int     `json:"rating"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}
