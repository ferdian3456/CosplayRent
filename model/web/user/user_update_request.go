package user

import "time"

type UserUpdateRequest struct {
	Id              string
	Name            string `json:"name"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	Password        string `json:"password"`
	Role            int    `json:"role"`
	Profile_picture string `json:"profile_picture"`
	Created_at      *time.Time
	Update_at       *time.Time
}
