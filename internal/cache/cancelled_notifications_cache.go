package cache

import (
	"context"
	"time"

	"github.com/wb-go/wbf/redis"
)

type ICancellationService interface {
	Cancel(ctx context.Context, notificationID string) error
	IsCancelled(ctx context.Context, notificationID string) (bool, error)
}

type RedisCancellationService struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedisCancellationService(client *redis.Client) ICancellationService {
	return &RedisCancellationService{
		client: client,
		prefix: "cancelled_notifications:",
		ttl:    24 * time.Hour,
	}
}

func (s *RedisCancellationService) key(notificationID string) string {
	return s.prefix + notificationID
}

func (s *RedisCancellationService) Cancel(ctx context.Context, notificationID string) error {
	err := s.client.SetWithExpiration(ctx, s.key(notificationID), "1", s.ttl)
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisCancellationService) IsCancelled(ctx context.Context, notificationID string) (bool, error) {
	result, err := s.client.Exists(ctx, s.key(notificationID)).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
