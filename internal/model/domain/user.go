package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id              uuid.UUID
	Name            string
	Email           string
	Address         *string
	Password        string
	Profile_picture *string
	Created_at      *time.Time
	Updated_at      *time.Time
}
