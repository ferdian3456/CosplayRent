package costume

import (
	"time"
)

type CostumeResponse struct {
	Id              int        `json:"id"`
	User_id         string     `json:"user_id"`
	Username        string     `json:"username"`
	Profile_picture *string    `json:"profile_picture"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Bahan           string     `json:"bahan"`
	Ukuran          string     `json:"ukuran"`
	Berat           string     `json:"berat"`
	Kategori        string     `json:"kategori"`
	Price           float64    `json:"price"`
	Picture         string     `json:"costume_picture"`
	Available       string     `json:"available"`
	Created_at      *time.Time `json:"created_at"`
	Updated_at      *time.Time `json:"updated_at"`
}
