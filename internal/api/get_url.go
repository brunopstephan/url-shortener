package api

import (
	"errors"
	"log/slog"
	"net/http"
	"urlShortener/internal/repositories"
	"urlShortener/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type getShortenedURLResponse struct {
	URL string `json:"url"`
}

func handleGetShortenedURL(db repositories.UrlRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		json := r.URL.Query().Get("json")

		data, err := db.GetURL(r.Context(), code)

		if err != nil {
			if errors.Is(err, redis.Nil) {
				utils.SendJSON(w, utils.ApiResponse{
					Error: "url not found",
				}, http.StatusNotFound)
				return
			}

			slog.Error("error get url", "error", err)
			utils.SendJSON(w, utils.ApiResponse{
				Error: "something went wrong",
			}, http.StatusInternalServerError)
			return
		}

		if json == "true" {
			utils.SendJSON(w, utils.ApiResponse{
				Data: getShortenedURLResponse{URL: data},
			}, http.StatusOK)
			return
		}
		http.Redirect(w, r, data, http.StatusMovedPermanently)

	}
}
