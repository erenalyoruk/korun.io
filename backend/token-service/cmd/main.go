package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"korun.io/token-service/internal/api"
	"korun.io/token-service/internal/config"
	"korun.io/token-service/internal/logging"
)

func main() {
	configPath := flag.String("config-path", "./configs", "Path to the configuration file")
	configName := flag.String("config-name", "config.dev", "Name of the configuration file without extension")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath, *configName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(cfg.Infra.Logger.Level, cfg.Infra.Logger.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	routes := api.SetupRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routes,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting server on port %d", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Server failed: %v", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.IdleTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server shutdown failed: %v", err))
	}

	logger.Info("Server gracefully stopped")
}
