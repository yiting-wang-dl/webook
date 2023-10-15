package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	Birthday time.Time // YYYY-MM-DD
	AboutMe  string

	// timezone，UTC 0 millisecond
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Address struct {
	Province string
	Region   string
}
