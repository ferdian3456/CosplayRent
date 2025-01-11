package domain

import "time"

type Costume struct {
	Id          int
	User_id     string
	Name        string
	Description string
	Bahan       string
	Ukuran      string
	Berat       int
	Kategori    int
	Price       float64
	Picture     string
	Available   string
	Created_at  *time.Time
	Updated_at  *time.Time
}
