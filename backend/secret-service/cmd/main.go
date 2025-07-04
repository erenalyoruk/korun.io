package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"korun.io/secret-service/internal/handlers"
	"korun.io/secret-service/internal/repository"
	"korun.io/secret-service/internal/service"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	repo := repository.NewSecretRepository(dbpool)
	svc := service.NewSecretService(repo)
	handler := handlers.NewSecretHandler(svc)

	r := gin.Default()

	r.GET("/secrets", handler.GetSecrets)
	r.GET("/secrets/:id", handler.GetSecretByID)
	r.GET("/secrets/name/:name", handler.GetSecretByName)
	r.POST("/secrets", handler.CreateSecret)
	r.PUT("/secrets/:id", handler.UpdateSecret)
	r.DELETE("/secrets/:id", handler.DeleteSecret)

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
