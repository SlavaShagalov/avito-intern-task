package main

import (
	"context"
	redisCache "github.com/SlavaShagalov/avito-intern-task/internal/banner/cache/redis"
	bannerRepository "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository/pgx"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	mw "github.com/SlavaShagalov/avito-intern-task/internal/middleware"
	pLog "github.com/SlavaShagalov/avito-intern-task/internal/pkg/log/zap"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/storage"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/storage/postgres"

	bannerDelivery "github.com/SlavaShagalov/avito-intern-task/internal/banner/delivery/http"
	bannerUsecase "github.com/SlavaShagalov/avito-intern-task/internal/banner/usecase"

	authDelivery "github.com/SlavaShagalov/avito-intern-task/internal/auth/delivery/http"
	authUsecase "github.com/SlavaShagalov/avito-intern-task/internal/auth/usecase"
	userRepository "github.com/SlavaShagalov/avito-intern-task/internal/user/repository/pgx"
)

func main() {
	// ===== Logger =====
	logger := pLog.NewDev()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()
	logger.Info("API server starting...")

	// ===== Configuration =====
	viper.SetConfigName("api")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Failed to read configuration", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("Configuration read successfully")

	// ===== Database =====
	pgxPool, err := postgres.NewPgx(logger)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		pgxPool.Close()
		logger.Info("Postgres connection closed")
	}()

	// ===== Cache =====
	ctx := context.Background()
	redisClient, err := storage.NewRedis(logger, ctx)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = redisClient.Close()
		if err != nil {
			logger.Error("Failed to close Redis connection", zap.Error(err))
		} else {
			logger.Info("Redis connection closed")
		}
	}()

	cache := redisCache.New(redisClient, logger)
	usersRepo := userRepository.New(pgxPool, logger)
	bannerRepo := bannerRepository.New(pgxPool, logger)

	authUC := authUsecase.New(usersRepo, logger)
	bannerUC := bannerUsecase.New(bannerRepo, logger)

	// ===== Server =====
	checkAuth := mw.NewCheckAuth(logger)
	checkAdminAccess := mw.NewCheckAdminAccess(logger)
	accessLog := mw.NewAccessLog(logger)
	panicCatch := mw.NewPanicCatch(logger)

	router := mux.NewRouter()

	authDelivery.RegisterHandlers(router, authUC, logger)
	bannerDelivery.RegisterHandlers(router, bannerUC, cache, logger, checkAuth, checkAdminAccess)

	server := http.Server{
		Addr:    ":" + viper.GetString(config.ServerPort),
		Handler: panicCatch(accessLog(router)),
	}

	logger.Info("API server started", zap.String("port", viper.GetString(config.ServerPort)))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped", zap.Error(err))
	}
}
