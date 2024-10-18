package user

import "time"

type UserCreateRequest struct {
	Name       string `validate:"required,min=5,max=20" json:"name"`
	Email      string `validate:"required,min=5,max=254" json:"email"`
	Password   string `validate:"required,min=5,max=20" json:"password"`
	Role       int    `validate:"required,min=1,max=1" json:"role"`
	Created_at *time.Time
}
