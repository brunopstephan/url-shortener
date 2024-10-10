package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	RedisHost     string
	RedisPort     string
	RedisPwd      string
	RedisDb       int
	BasicAuthUser string
	BasicAuthPwd  string
}

func getConfig() config {
	err := godotenv.Load()

	if err != nil {
		slog.Error("Error loading .env file")
		panic(err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPwd := os.Getenv("REDIS_PASSWORD")
	basicAuthUser := os.Getenv("BASIC_AUTH_USER")
	basicAuthPwd := os.Getenv("BASIC_AUTH_PASSWORD")
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		slog.Error("error converting redis db to int", "error", err)
		panic(err)
	}

	return config{
		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisPwd:      redisPwd,
		RedisDb:       redisDb,
		BasicAuthUser: basicAuthUser,
		BasicAuthPwd:  basicAuthPwd,
	}
}

var Config = getConfig()
