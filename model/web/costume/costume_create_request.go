package costume

import (
	"time"
)

type CostumeCreateRequest struct {
	User_id     string  `validate:"required,min=36,max=36" json:"user_id"`
	Name        string  `validate:"required,min=5,max=30" json:"name"`
	Description string  `validate:"required,min=5,max=1000" json:"description"`
	Bahan       string  `validate:"required,min=3,max=30" json:"bahan"`
	Ukuran      string  `validate:"required,min=3,max=30" json:"ukuran"`
	Berat       string  `validate:"required,min=3,max=30" json:"berat"`
	Kategori    string  `validate:"required,min=3,max=30" json:"kategori"`
	Price       float64 `validate:"required" json:"price"`
	Picture     string  `validate:"required" json:"picture"`
	Available   bool    `json:"available"`
	Created_at  *time.Time
}
