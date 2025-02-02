package costume

type CostumeUpdateRequest struct {
	Id          int
	Name        string  `validate:"max=30" json:"name,omitempty"`
	Description string  `validate:"max=1000" json:"description,omitempty"`
	Bahan       string  `validate:"max=30" json:"bahan,omitempty"`
	Ukuran      string  `validate:"max=30" json:"ukuran,omitempty"`
	Berat       int     `json:"berat,omitempty"`
	Kategori    int     `validate:"max=30" json:"kategori,omitempty"`
	Available   string  `validate:"max=13" json:"available,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Picture     *string `json:"costume_picture,omitempty"`
}
