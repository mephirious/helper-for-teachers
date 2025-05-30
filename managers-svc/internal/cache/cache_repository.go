package cache

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/ports"
)

type CacheRepository struct {
	repo  ports.MemberRepository
	cache map[string]interface{}
	mu    sync.RWMutex
}

func NewCacheRepository(repo ports.MemberRepository) *CacheRepository {
	c := &CacheRepository{
		repo:  repo,
		cache: make(map[string]interface{}),
	}

	return c
}

func (c *CacheRepository) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupMember, error) {
	c.mu.RLock()
	if members, exists := c.cache[groupID.String()+"members"]; exists {
		c.mu.RUnlock()
		return members.([]*domain.GroupMember), nil
	}
	c.mu.RUnlock()

	members, err := c.repo.ListByGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[groupID.String()+"members"] = members
	c.mu.Unlock()

	return members, nil
}

func (c *CacheRepository) Create(ctx context.Context, member *domain.GroupMember) error {
	if err := c.repo.Create(ctx, member); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	members, err := c.repo.ListByGroup(ctx, member.GroupID)
	if err != nil {
		return err
	}
	c.cache[member.GroupID.String()+"members"] = members

	return nil
}

func (c *CacheRepository) Delete(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) error {
	if err := c.repo.Delete(ctx, groupID, userID, role); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	members, err := c.repo.ListByGroup(ctx, groupID)
	if err != nil {
		return err
	}
	c.cache[groupID.String()+"members"] = members

	return nil
}

func (c *CacheRepository) Exists(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) (bool, error) {
	return c.repo.Exists(ctx, groupID, userID, role)
}
