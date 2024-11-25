package web

import "time"

type UserCreateRequest struct {
	Name       string `validate:"required,min=5,max=20" json:"name"`
	Email      string `validate:"required,min=5,max=254" json:"email"`
	Password   string `validate:"required,min=5,max=20" json:"password"`
	Created_at *time.Time
}

type UserResponse struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Address         *string `json:"address"`
	Profile_picture *string `json:"profile_picture"`
	Created_at      string  `json:"created_at"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,min=5,max=254" json:"email"`
	Password string `validate:"required,min=5,max=20" json:"password"`
}

type UserUpdateRequest struct {
	Id              string
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Address         *string `json:"address"`
	Profile_picture *string `json:"profile_picture"`
	Created_at      *time.Time
	Update_at       *time.Time
}

type FixUser struct {
	Id string
}
