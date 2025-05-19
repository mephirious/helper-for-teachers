package messaging

import (
	"context"

	eventspb "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
)

type EventPublisher interface {
	PublishEventCreated(ctx context.Context, event *eventspb.Event) error
	PublishEventUpdated(ctx context.Context, event *eventspb.Event) error
	PublishEventDeleted(ctx context.Context, eventID string) error
	Close()
}
