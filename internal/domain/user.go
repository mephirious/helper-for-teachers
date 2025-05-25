package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `bson:"_id"`
	Email     string    `bson:"email"`
	Username  string    `bson:"username"`
	Password  string    `bson:"password"`
	Role      Role      `bson:"role"`
	Phone     string    `bson:"phone"`
	Verified  bool      `bson:"verified"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UpdateUserProfileParams struct {
	ID       string
	Email    string
	Username string
	Phone    string
}

type Role string

const (
	UNSPECIFIED Role = ""
	ADMIN       Role = "admin"
	TEACHER     Role = "teacher"
	STUDENT     Role = "student"
)

func NewUser(email, hashedPwd string, role Role) *User {
	return &User{
		ID:        uuid.NewString(),
		Email:     email,
		Password:  hashedPwd,
		Role:      role,
		Verified:  false,
		CreatedAt: time.Now().UTC(),
	}
}

