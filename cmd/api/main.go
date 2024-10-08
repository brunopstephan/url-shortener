package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
	"urlShortener/internal/api"
	"urlShortener/internal/repositories"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed initializing the application", "error", err)
		return
	}
	slog.Info("All systems offline")
}

func run() error {
	err := godotenv.Load()

	if err != nil {
		slog.Error("Error loading .env file")
		return err
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	redisAddr := redisHost + ":" + redisPort

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	urlRepository := repositories.NewUrlRepository(rdb)
	handler := api.NewHandler(urlRepository)
	s := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
		Addr:         ":9000",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
