package api

import (
	"net/http"
	"strconv"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"

	_ "url-shortener/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(db repositories.UrlContract) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	port := config.Config.Port
	_url := "http://localhost:" + strconv.Itoa(port) + "/swagger/doc.json"

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(_url), //The url pointing to API definition
	))

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlers.HandlePostShortenedURL(db))
		r.Get("/{code}", handlers.HandleGetShortenedURL(db))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("Restricted", map[string]string{
			config.Config.BasicAuthUser: config.Config.BasicAuthPwd,
		}))

		r.Route("/admin", func(r chi.Router) {
			r.Get("/all", handlers.HandleGetAllUrls(db))
			r.Delete("/{code}", handlers.HandleDeleteShortenedURL(db))
			r.Put("/{code}", handlers.HandleUpdateShortenedURL(db))
		})
	})
	return r
}
