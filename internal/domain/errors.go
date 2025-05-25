package domain

import "errors"

var (
	// User errors
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrPasswordUnchanged     = errors.New("password unchanged")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidArgument       = errors.New("must specify either ID or Email")
	ErrPermissionDenied      = errors.New("permission denied")

	// Token errors
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")

	// Code errors
	ErrCodeInvalid              = errors.New("verification code invalid")
	ErrCodeExpired              = errors.New("verification code expired")
	ErrPasswordResetCodeExpired = errors.New("password reset code expired")
	ErrPasswordResetCodeInvalid = errors.New("invalid password reset code")
	ErrInvalidPurpose           = errors.New("invalid code purpose")
)
