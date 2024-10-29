package user

type UserLoginRequest struct {
	Email    string `validate:"required,min=5,max=20" json:"email"`
	Password string `validate:"required,min=5,max=20" json:"password"`
}
