package user

import "time"

type UserCreateRequest struct {
	Name       string `validate:"required,min=5,max=20" json:"name"`
	Email      string `validate:"required,min=5,max=254" json:"email"`
	Password   string `validate:"required,min=5,max=20" json:"password"`
	Created_at *time.Time
}

type IdentityCardRequest struct {
	IdentityCard_picture string `json:"identitycard_picture"`
}
