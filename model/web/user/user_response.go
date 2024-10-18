package user

import (
	"time"
)

type UserResponse struct {
	Id              string     `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Address         *string    `json:"address"`
	Role            int        `json:"role"`
	Profile_picture *string    `json:"profile_picture"`
	Created_at      *time.Time `json:"created_at"`
}
