package models

import "time"

type User struct {
	ID        int64
	Username  string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
}
