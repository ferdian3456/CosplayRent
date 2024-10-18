package costume

import (
	"time"
)

type CostumeResponse struct {
	Id          int        `json:"id"`
	User_id     string     `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Picture     *string    `json:"picture"`
	Available   string     `json:"available"`
	Created_at  *time.Time `json:"created_at"`
}
