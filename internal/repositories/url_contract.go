package repositories

import "context"

type UrlContract interface {
	SaveShortenedURL(ctx context.Context, _url string) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
	GetAllURL(ctx context.Context) (map[string]string, error)
}
