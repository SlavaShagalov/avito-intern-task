package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type Cache struct {
	rdb *redis.Client
	log *zap.Logger
}

func New(rdb *redis.Client, log *zap.Logger) *Cache {
	return &Cache{
		rdb: rdb,
		log: log,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	err := c.rdb.Set(ctx, key, value, 5*time.Minute).Err()
	if err != nil {
		c.log.Error("Failed to set key-value in Redis", zap.Error(err))
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	var value []byte
	err := c.rdb.Get(ctx, key).Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
