package ports

type GRPCPort interface {
	eventsv1.EventServiceServer
}
