package costume

import (
	"time"
)

type CostumeCreateRequest struct {
	User_id     string
	Name        string  `validate:"required,min=5,max=30" json:"name"`
	Description string  `validate:"required,min=5,max=1000" json:"description"`
	Bahan       string  `validate:"required,min=3,max=30" json:"bahan"`
	Ukuran      string  `validate:"required,min=3,max=30" json:"ukuran"`
	Berat       int     `validate:"required,min=1" json:"berat"`
	Kategori    string  `validate:"required,min=3,max=30" json:"kategori"`
	Price       float64 `validate:"required" json:"price"`
	Picture     *string `json:"costume_picture"`
	Available   bool    `json:"available"`
	Created_at  *time.Time
}
