package domain

import "time"

type UserVerification struct {
	Id                int
	User_id           string
	Verification_code string
	Created_at        *time.Time
	Updated_at        *time.Time
	Expired_at        *time.Time
}
