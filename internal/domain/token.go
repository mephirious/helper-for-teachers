package domain

import (
	"time"
)

type TokenPayload struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`  // UTC
	ExpiresAt int64     `json:"expires_at"` // Unix
}

type VerificationCode struct {
	UserID    string
	Code      string
	ExpiresAt time.Time
	Purpose   string
}

// Purposes for code verification
const (
	PurposeEmailVerification = "email_verification"
	PurposeResetPassword     = "reset_password"
)

// Creates a verification code with TTL
func NewVerificationCode(userID, code, purpose string, ttl time.Duration) *VerificationCode {
	return &VerificationCode{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(ttl),
	}
}
