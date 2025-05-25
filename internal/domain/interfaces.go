package domain

import (
	"context"
	"time"
)

// Helper services
type PasswordHasher interface {
	Hash(ctx context.Context, plain string) (string, error)
	Verify(ctx context.Context, hashed, plain string) bool
}

type EmailSender interface {
	Send(to, subject, htmlBody string) error
}

type JWTService interface {
	Generate(userID string, role Role) (accessToken string, issuedAt time.Time, expiresAt int64, err error) // generate access_token
	Validate(ctx context.Context, token string) (*TokenPayload, error)                                      // validate access_token
}

type CodeCache interface {
	Set(ctx context.Context, code *VerificationCode) error
	Get(ctx context.Context, userID string) (*VerificationCode, error)
	Delete(ctx context.Context, userID string) error
}
