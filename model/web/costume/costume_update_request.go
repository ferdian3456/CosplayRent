package costume

import "time"

type CostumeUpdateRequest struct {
	Id          int
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Bahan       *string  `json:"bahan"`
	Ukuran      *string  `json:"ukuran"`
	Berat       *string  `json:"berat"`
	Kategori    *string  `json:"kategori"`
	Price       *float64 `json:"price"`
	Picture     *string  `json:"picture"`
	Available   bool     `json:"available"`
	Created_at  time.Time
	Update_at   *time.Time
}
