package ports

import (
	"context"

	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type NatsPublisher interface {
	PublishCourseCreated(ctx context.Context, course *domain.Course) error
	PublishGroupCreated(ctx context.Context, group *domain.Group) error
	PublishGroupMemberAdded(ctx context.Context, member *domain.GroupMember) error
	PublishCourseInstructorAssigned(ctx context.Context, instr *domain.CourseInstructor) error
}
