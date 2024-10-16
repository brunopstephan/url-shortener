package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"urlShortener/internal/api"
	"urlShortener/internal/config"
	"urlShortener/internal/repositories"

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
	redisAddr := config.Config.RedisHost + ":" + config.Config.RedisPort

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: config.Config.RedisPwd,
		DB:       config.Config.RedisDb,
	})

	urlRepository := repositories.NewUrlRepository(rdb)
	handler := api.NewHandler(urlRepository)
	s := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
		Addr:         ":" + strconv.Itoa(config.Config.Port),
		Handler:      handler,
	}

	slog.Info(fmt.Sprintf("Server started on port %d", config.Config.Port))
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
