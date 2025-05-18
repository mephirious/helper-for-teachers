package ports

import eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"

type GRPCPort interface {
	eventsv1.EventServiceServer
}
