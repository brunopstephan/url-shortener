package api

import (
	"net/http"
	"urlShortener/internal/handlers"
	"urlShortener/internal/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db repositories.UrlContract) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlers.HandlePostShortenedURL(db))
		r.Get("/{code}", handlers.HandleGetShortenedURL(db))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("Restricted", map[string]string{
			"admin": "password",
		}))

		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/all", handlers.HandleGetAllUrls(db))
		})
	})
	return r
}
