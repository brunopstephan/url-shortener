package api

import (
	"net/http"
	"urlShortener/internal/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db repositories.UrlRepositoryInterface) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlePostShortenedURL(db))
		r.Get("/{code}", handleGetShortenedURL(db))

	})
	return r
}
