package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"url-shortener/internal/api"
	"url-shortener/internal/config"
	"url-shortener/internal/repositories"

	"github.com/redis/go-redis/v9"
)

// @title URL Shortener API
// @version 1.0
// @description A simple url shortener.

// @contact.name Bruno Piffer Stephan
// @contact.email brunopstephan@gmail.com

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
		Addr:         ":" + strconv.Itoa(config.Config.AppPort),
		Handler:      handler,
	}

	slog.Info(fmt.Sprintf("Server started on port %d", config.Config.AppPort))
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
