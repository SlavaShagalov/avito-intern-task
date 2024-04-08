package postgres

import (
	"context"
	"fmt"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewPgx(log *zap.Logger) (*pgxpool.Pool, error) {
	log.Info("Postgres PGX connecting...",
		zap.String("host", viper.GetString(config.PostgresHost)),
		zap.Int("port", viper.GetInt(config.PostgresPort)),
		zap.String("db", viper.GetString(config.PostgresDB)),
	)

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=100",
		viper.GetString(config.PostgresUser),
		viper.GetString(config.PostgresPassword),
		viper.GetString(config.PostgresHost),
		strconv.Itoa(viper.GetInt(config.PostgresPort)),
		viper.GetString(config.PostgresDB),
	)

	conf, _ := pgxpool.ParseConfig(connString)
	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		log.Error("Failed to connect to Postgres PGX", zap.Error(err))
		return nil, err
	}

	log.Info("Postgres PGX connected")
	return pool, nil
}
