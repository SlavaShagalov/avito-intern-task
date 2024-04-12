package storage

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewRedis(log *zap.Logger, ctx context.Context) (*redis.Client, error) {
	log.Info("Redis connecting...",
		zap.String("host", viper.GetString(config.RedisHost)),
		zap.String("port", viper.GetString(config.RedisPort)),
		zap.Int("db", viper.GetInt(config.RedisDB)),
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(config.RedisHost) + ":" + viper.GetString(config.RedisPort),
		Password: viper.GetString(config.RedisPassword),
		DB:       viper.GetInt(config.RedisDB),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to create Redis connection, ", zap.Error(err))
		return nil, err
	}

	log.Info("Redis connected")
	return rdb, nil
}
