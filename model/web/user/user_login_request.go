package user

type UserLoginRequest struct {
	Name     string `validate:"required,min=5,max=20" json:"name"`
	Password string `validate:"required,min=5,max=20" json:"password"`
}
