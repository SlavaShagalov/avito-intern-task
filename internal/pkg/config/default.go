package config

import (
	"github.com/spf13/viper"
)

// Postgres

func SetTestPostgresConfig() {
	viper.SetDefault(PostgresHost, "localhost")
	viper.SetDefault(PostgresPort, 5435)
	viper.SetDefault(PostgresDB, "banners_db")
	viper.SetDefault(PostgresUser, "moderator")
	viper.SetDefault(PostgresPassword, "2222")
	viper.SetDefault(PostgresSSLMode, "disable")
}

// Redis

func SetTestRedisConfig() {
	viper.SetDefault(RedisHost, "localhost")
	viper.SetDefault(RedisPort, 6379)
	viper.SetDefault(RedisPassword, "2222")
	viper.SetDefault(RedisDB, 0)
}
