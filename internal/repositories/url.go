package repositories

import (
	"context"
	"errors"
	"fmt"
	"urlShortener/internal/utils"

	"github.com/redis/go-redis/v9"
)

type UrlRepository struct {
	rdb *redis.Client
}

type UrlRepositoryInterface interface {
	SaveShortenedURL(ctx context.Context, _url string) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}

func NewUrlRepository(rdb *redis.Client) UrlRepositoryInterface {
	return &UrlRepository{rdb: rdb}
}

func (s *UrlRepository) SaveShortenedURL(ctx context.Context, _url string) (string, error) {
	var code string
	for range 5 {
		code = utils.GenCode()
		if err := s.rdb.HGet(ctx, "encurtador", code).Err(); err != nil {
			if errors.Is(err, redis.Nil) {
				break
			}
			return "", fmt.Errorf("code already exists: %w", err)
		}
	}

	if err := s.rdb.HSet(ctx, "encurtador", code, _url).Err(); err != nil {
		return "", fmt.Errorf("error setting on redis: %w", err)
	}

	return code, nil
}

func (s *UrlRepository) GetURL(ctx context.Context, code string) (string, error) {
	_url, err := s.rdb.HGet(ctx, "encurtador", code).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	return _url, nil
}