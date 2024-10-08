package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"urlShortener/internal/repositories"
	"urlShortener/internal/utils"
)

type postBody struct {
	URL string `json:"url"`
}

func handlePostShortenedURL(db repositories.UrlRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body postBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.SendJSON(w, utils.ApiResponse{Error: "invalid request body"}, http.StatusUnprocessableEntity)
			return
		}

		if body.URL == "" {
			utils.SendJSON(w, utils.ApiResponse{Error: "URL is required"}, http.StatusBadRequest)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			utils.SendJSON(w, utils.ApiResponse{Error: "invalid URL"}, http.StatusBadRequest)
			return
		}

		code, err := db.SaveShortenedURL(r.Context(), body.URL)
		if err != nil {
			slog.Error("error saving url", "error", err)
			utils.SendJSON(w, utils.ApiResponse{
				Error: "something went wrong",
			}, http.StatusInternalServerError)
			return
		}

		utils.SendJSON(w, utils.ApiResponse{Data: code}, http.StatusCreated)

	}
}
