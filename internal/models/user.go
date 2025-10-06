package models

import "time"

type User struct {
	ID            int
	Username      string
	FirstName     string
	LastName      string
	Email         string
	Password      string
	IsVerified    bool
	IsAdmin       bool
	AccessLevel   int
	SignupIP      string
	SignupCountry string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
