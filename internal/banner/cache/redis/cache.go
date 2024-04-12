package redis

import (
	"context"
	"encoding/json"
	"github.com/SlavaShagalov/avito-intern-task/internal/banner/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

const (
	expiration = 5 * time.Minute
)

type redisCache struct {
	rdb *redis.Client
	log *zap.Logger
}

func New(rdb *redis.Client, log *zap.Logger) cache.Cache {
	return &redisCache{
		rdb: rdb,
		log: log,
	}
}

func (c *redisCache) Set(ctx context.Context, key string, value *cache.Value) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		c.log.Error("Cache: failed to marshal value", zap.Error(err))
		return err
	}
	err = c.rdb.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		c.log.Error("Cache: failed to set key-value", zap.Error(err))
		return err
	}
	return nil
}

func (c *redisCache) Get(ctx context.Context, key string) (*cache.Value, error) {
	var jsonValue []byte
	err := c.rdb.Get(ctx, key).Scan(&jsonValue)
	if err != nil {
		c.log.Debug("Cache: failed to get value", zap.Error(err))
		return nil, err
	}
	value := new(cache.Value)
	err = json.Unmarshal(jsonValue, value)
	if err != nil {
		c.log.Error("Cache: failed to unmarshal value", zap.Error(err))
		return nil, err
	}
	return value, nil
}
