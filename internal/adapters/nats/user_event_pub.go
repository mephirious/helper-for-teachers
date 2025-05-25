package nats

import (
	"context"

	"github.com/mephirious/helper-for-teachers/services/auth-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/auth-svc/pkg/nats"
)

type AuthPublisher struct {
	conn *nats.Client
}

func NewAuthPublisher(c *nats.Client) *AuthPublisher {
	return &AuthPublisher{conn: c}
}

func (p *AuthPublisher) PublishUserRegistered(ctx context.Context, evt *domain.UserRegisteredEvent) error {
	return p.conn.PublishJSON("user.registered", evt)
}
