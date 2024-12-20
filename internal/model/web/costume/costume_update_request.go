package costume

type CostumeUpdateRequest struct {
	Id          int
	Name        string  `validate:"min=5,max=30" json:"name,omitempty"`
	Description string  `validate:"min=5,max=1000" json:"description,omitempty"`
	Bahan       string  `validate:"min=5,max=30" json:"bahan,omitempty"`
	Ukuran      string  `validate:"min=1,max=30" json:"ukuran,omitempty"`
	Berat       int     `json:"berat,omitempty"`
	Kategori    string  `validate:"min=3,max=30" json:"kategori,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Picture     *string `validate:"max=254" json:"costume_picture,omitempty"`
}
