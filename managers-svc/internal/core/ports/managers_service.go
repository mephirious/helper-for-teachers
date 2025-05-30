package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type ManagersService interface {
	CourseService
	GroupService
	MemberService
	InstructorService
}

type CourseService interface {
	CreateCourse(ctx context.Context, name string) (*domain.Course, error)
	GetCourse(ctx context.Context, id uuid.UUID) (*domain.Course, error)
	UpdateCourse(ctx context.Context, update domain.CourseUpdate) (*domain.Course, error)
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	ListCourses(ctx context.Context) ([]*domain.Course, error)
}

type GroupService interface {
	CreateGroup(ctx context.Context, courseID uuid.UUID, name string) (*domain.Group, error)
	GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error)
	UpdateGroup(ctx context.Context, update domain.GroupUpdate) (*domain.Group, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	ListGroupsByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.Group, error)
}

type MemberService interface {
	AddGroupMember(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) (*domain.GroupMember, error)
	RemoveGroupMember(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) error
	ListGroupMembers(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupMember, error)
	IsUserGroupTeacher(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
}

type InstructorService interface {
	AssignCourseInstructor(ctx context.Context, courseID, userID uuid.UUID) (*domain.CourseInstructor, error)
	RemoveCourseInstructor(ctx context.Context, courseID, userID uuid.UUID) error
	ListCourseInstructors(ctx context.Context, courseID uuid.UUID) ([]*domain.CourseInstructor, error)
}
