package domain

import "time"

type Notification struct {
	Id               int
	Template_name    string
	Template_subject string
	Template_body    string
	Created_at       *time.Time
	Updated_at       *time.Time
}

type EmailNotification struct {
	Username string
	Code     string
}
