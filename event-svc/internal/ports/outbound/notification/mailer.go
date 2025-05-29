package notification

import (
	"context"
	"time"

	eventspb "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
)

type Mailer interface {
	SendLessonNotification(ctx context.Context, lesson *eventspb.Lesson, recipient *eventspb.User) error
	SendTaskNotification(ctx context.Context, task *eventspb.Task, recipient *eventspb.User) error
	ScheduleNotifications(ctx context.Context, scheduler SchedulerService)
	Close()
}

type SchedulerService interface {
	GetUpcomingLessons(ctx context.Context, from, to time.Time) ([]*eventspb.Lesson, error)
	GetUpcomingTasks(ctx context.Context, from, to time.Time) ([]*eventspb.Task, error)
}
