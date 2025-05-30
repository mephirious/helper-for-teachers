package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"updated_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Group struct {
	ID        uuid.UUID `db:"id"`
	CourseID  uuid.UUID `db:"course_id"`
	Name      string    `db:"name"`
	ExpireAt  time.Time `db:"expire_at"`
	CreatedAt time.Time `db:"updated_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type MemberRole string

const (
	StudentRole MemberRole = "student"
	TeacherRole MemberRole = "teacher"
)

func NewMemberRoleFromProto(protoRole string) (MemberRole, error) {
	switch protoRole {
	case "STUDENT":
		return StudentRole, nil
	case "TEACHER":
		return TeacherRole, nil
	default:
		return "", fmt.Errorf("invalid role: %s", protoRole)
	}
}

type GroupMember struct {
	ID        uuid.UUID  `db:"id"`
	GroupID   uuid.UUID  `db:"group_id"`
	UserID    uuid.UUID  `db:"user_id"`
	Role      MemberRole `db:"role"`
	CreatedAt time.Time  `db:"updated_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

type CourseInstructor struct {
	ID        uuid.UUID `db:"id"`
	CourseID  uuid.UUID `db:"course_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"updated_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CourseUpdate struct {
	ID   uuid.UUID
	Name *string
}
type GroupUpdate struct {
	ID       uuid.UUID
	Name     *string
	ExpireAt *time.Time
}

var (
	ErrInvalidRequestPayload    = fmt.Errorf("invalid request payload")
	ErrNotImplemented           = fmt.Errorf("not implemented")
	ErrNotUpdated               = fmt.Errorf("not updated")
	ErrGroupNotFound            = fmt.Errorf("group not found")
	ErrCourseNotFound           = fmt.Errorf("course not found")
	ErrCourseInstructorNotFound = fmt.Errorf("course instructor not found")
	ErrGroupMemberNotFound      = fmt.Errorf("group member not found")
)
