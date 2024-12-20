package domain

import "time"

type TopUpOrder struct {
	Id             string
	User_id        string
	TopUp_amount   float64
	Status_payment bool
	Created_at     *time.Time
	Updated_at     *time.Time
}
