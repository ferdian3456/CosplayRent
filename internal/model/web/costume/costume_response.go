package costume

type CostumeResponse struct {
	Id              int     `json:"id"`
	User_id         string  `json:"user_id"`
	Username        string  `json:"username"`
	Profile_picture *string `json:"profile_picture"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Bahan           string  `json:"bahan"`
	Ukuran          *string `json:"ukuran"`
	Berat           int     `json:"berat"`
	Kategori        string  `json:"kategori"`
	Kategori_id     int     `json:"-"`
	Price           float64 `json:"price"`
	Picture         *string `json:"costume_picture"`
	Available       string  `json:"available"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}

type SimpleCostumeResponse struct {
	Name            string `json:"name"`
	Costume_picture string `json:"costume_picture"`
	Costume_size    string `json:"size"`
	Costume_weight  int    `json:"weight"`
}

type SellerCostumeResponse struct {
	Id          int     `json:"id"`
	User_id     string  `json:"user_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Bahan       string  `json:"bahan"`
	Ukuran      *string `json:"ukuran"`
	Berat       int     `json:"berat"`
	Kategori    string  `json:"kategori"`
	Price       float64 `json:"price"`
	KotaAsal    string  `json:"kota_asal"`
	Picture     *string `json:"costume_picture"`
	Available   string  `json:"available"`
	Created_at  string  `json:"created_at"`
	Updated_at  string  `json:"updated_at"`
}
