package costume

type CostumeCreateRequest struct {
	Name        string  `validate:"required,min=5,max=30" json:"name"`
	Description string  `validate:"required,min=5,max=1000" json:"description"`
	Bahan       string  `validate:"required,min=3,max=30" json:"bahan"`
	Ukuran      string  `validate:"required,min=3,max=30" json:"ukuran"`
	Berat       int     `validate:"required,min=1" json:"berat"`
	Kategori    string  `validate:"required,min=3,max=30" json:"kategori"`
	Price       float64 `validate:"required" json:"price"`
	Picture     *string `validate:"required,max=254" json:"costume_picture"`
}
