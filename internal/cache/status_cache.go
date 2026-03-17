package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"WBTech_L3.1/internal/model"
	"github.com/wb-go/wbf/redis"
)

type StatusCache interface {
	Set(ctx context.Context, n *model.Notification) error
	Get(ctx context.Context, id string) (*model.Notification, error)
	Delete(ctx context.Context, id string) error
}

type redisStatusCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewStatusCache(client *redis.Client) StatusCache {
	if client == nil {
		return &noopStatusCache{}
	}
	return &redisStatusCache{client: client, ttl: 10 * time.Minute}
}

func (c *redisStatusCache) Set(ctx context.Context, n *model.Notification) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return c.client.SetWithExpiration(ctx, cacheKey(n.ID), b, c.ttl)
}

func (c *redisStatusCache) Get(ctx context.Context, id string) (*model.Notification, error) {
	res, err := c.client.Get(ctx, cacheKey(id))
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			return nil, nil
		}
		return nil, err
	}
	var n model.Notification
	if err = json.Unmarshal([]byte(res), &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func (c *redisStatusCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, cacheKey(id))
}

type noopStatusCache struct{}

func (n *noopStatusCache) Set(ctx context.Context, nn *model.Notification) error { return nil }
func (n *noopStatusCache) Get(ctx context.Context, id string) (*model.Notification, error) {
	return nil, nil
}
func (n *noopStatusCache) Delete(ctx context.Context, id string) error { return nil }

func cacheKey(id string) string {
	return "notify:" + id
}
