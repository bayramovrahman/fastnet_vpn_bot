package models

import "time"

type UserLoginSecurity struct {
	ID                     int
	UserID                 int
	EmailVerification      bool
	PhoneVerification      bool
	MultiFactorAuth        bool
	VerificationCode       string
	CodeExpiresAt          time.Time
	PhoneNumber            string
	LastVerificationSentAt time.Time
	FailedAttempts         int
	LockedUntil            time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
