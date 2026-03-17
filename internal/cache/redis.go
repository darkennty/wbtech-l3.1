package cache

import (
	"WBTech_L3.1/internal/config"
	"github.com/wb-go/wbf/redis"
)

func NewRedisClient(cfg config.Config) *redis.Client {
	if cfg.RedisAddr == "" {
		return nil
	}
	return redis.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
}
