package domain

import "time"

type OrderEvents struct {
	Id                       int
	User_id                  string
	Order_id                 string
	Status                   string
	Notes                    string
	Shipment_receipt_user_id string
	Created_at               *time.Time
}
