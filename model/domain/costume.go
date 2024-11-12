package domain

import "time"

type Costume struct {
	Id          int
	User_id     string
	Name        string
	Description string
	Bahan       string
	Ukuran      string
	Berat       string
	Kategori    string
	Price       float64
	Picture     string
	Available   bool
	Created_at  *time.Time
	Update_at   *time.Time
}
