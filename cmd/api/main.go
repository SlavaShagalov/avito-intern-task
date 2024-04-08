package main

import (
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	mw "github.com/SlavaShagalov/avito-intern-task/internal/middleware"
	pLog "github.com/SlavaShagalov/avito-intern-task/internal/pkg/log/zap"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/storage/postgres"

	bannerDelivery "github.com/SlavaShagalov/avito-intern-task/internal/banner/delivery/http"
	bannerRepository "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository/pgx"
	bannerUsecase "github.com/SlavaShagalov/avito-intern-task/internal/banner/usecase"
)

func main() {
	// ===== Logger =====
	logger := pLog.NewDevelop()
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
	viper.AddConfigPath("/configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read configuration: %v\n", err)
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

	bannerRepo := bannerRepository.New(pgxPool, logger)
	bannerUC := bannerUsecase.New(bannerRepo, logger)

	// ===== Server =====
	accessLog := mw.NewAccessLog(logger)
	router := mux.NewRouter()

	bannerDelivery.RegisterHandlers(router, bannerUC, logger)

	server := http.Server{
		Addr:    ":" + viper.GetString(config.ServerPort),
		Handler: accessLog(router),
	}

	logger.Info("API server started", zap.String("port", viper.GetString(config.ServerPort)))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped", zap.Error(err))
	}
}
