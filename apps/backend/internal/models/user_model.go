package models

import "time"

type User struct {
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	LastLogin    *time.Time
}
