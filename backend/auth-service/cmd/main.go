package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"korun.io/auth-service/internal/api"
	"korun.io/auth-service/internal/config"
	"korun.io/auth-service/internal/infrastructure/redis"
	"korun.io/auth-service/internal/infrastructure/persistence"
	"korun.io/auth-service/internal/application"
	"korun.io/shared/messaging"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	setupLogger(cfg.Infra.Logger.Level, cfg.Infra.Logger.Format)

	slog.Info("Starting Auth Service", "version", "1.0.0")

	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := redis.NewClient(&cfg.Infra.Redis)
	if err != nil {
		slog.Error("Failed to initialize Redis client", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	producer, err := messaging.NewKafkaProducer(&cfg.Infra.Kafka)
	if err != nil {
		slog.Error("Failed to initialize Kafka producer", "error", err)
		os.Exit(1)
	}
	defer producer.Close()

	accountRepo := persistence.NewPostgresAuthRepository(db)
	tokenRepo := persistence.NewRedisRefreshTokenRepository(redisClient)

	tokenService := application.NewTokenService(tokenRepo, &cfg.JWT)
	authService := application.NewAuthService(accountRepo, tokenService, producer, &cfg.Infra)

	router := api.SetupRoutes(authService, &cfg.Server)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("Server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}

func setupLogger(level, format string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: logLevel}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
