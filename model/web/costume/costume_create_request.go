package costume

import (
	"time"
)

type CostumeCreateRequest struct {
	User_id     string  `validate:"required,min=36,max=36" json:"user_id"`
	Name        string  `validate:"required,min=5,max=30" json:"name"`
	Description string  `validate:"required,min=5,max=1000" json:"description"`
	Price       float64 `validate:"required,min=1,max=1000" json:"price"`
	Picture     string  `validate:"required,min=1,max=100000" json:"picture"`
	Available   bool    `json:"available"`
	Created_at  *time.Time
}
