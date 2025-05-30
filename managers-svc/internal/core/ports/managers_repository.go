package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type CourseRepository interface {
	Create(ctx context.Context, course *domain.Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Course, error)
	Update(ctx context.Context, id uuid.UUID, updateFn func(*domain.Course) (bool, error)) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*domain.Course, error)
}

type GroupRepository interface {
	Create(ctx context.Context, group *domain.Group) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Group, error)
	Update(ctx context.Context, id uuid.UUID, updateFn func(*domain.Group) (bool, error)) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.Group, error)
}

type InstructorRepository interface {
	Create(ctx context.Context, instr *domain.CourseInstructor) error
	Delete(ctx context.Context, courseID, userID uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.CourseInstructor, error)
}

type MemberRepository interface {
	Create(ctx context.Context, gm *domain.GroupMember) error
	Delete(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) error
	ListByGroup(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupMember, error)
	Exists(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) (bool, error)
}
