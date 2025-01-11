package user

type UserCreateRequest struct {
	Name     string `validate:"required,min=5,max=20" json:"name"`
	Email    string `validate:"required,min=5,max=254" json:"email"`
	Password string `validate:"required,min=5,max=20" json:"password"`
}

type IdentityCardRequest struct {
	IdentityCard_picture *string `validate:"required,min=5,max=255" json:"identitycard_picture"`
}
