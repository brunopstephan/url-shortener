package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"urlShortener/internal/repositories"
	"urlShortener/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type getShortenedURLResponse struct {
	URL string `json:"url"`
}

func HandleGetShortenedURL(db repositories.UrlContract) http.HandlerFunc {
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

type postBody struct {
	URL string `json:"url"`
}

func HandlePostShortenedURL(db repositories.UrlContract) http.HandlerFunc {
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

type getAllUrlsResponse struct {
	URLs map[string]string `json:"urls"`
}

func HandleGetAllUrls(db repositories.UrlContract) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := db.GetAllURL(r.Context())
		if err != nil {
			slog.Error("error get urls", "error", err)
			utils.SendJSON(w, utils.ApiResponse{
				Error: "something went wrong",
			}, http.StatusInternalServerError)
			return
		}

		utils.SendJSON(w, utils.ApiResponse{
			Data: getAllUrlsResponse{URLs: urls},
		}, http.StatusOK)

	}
}

func HandleDeleteShortenedURL(db repositories.UrlContract) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		if err := db.DeleteURL(r.Context(), code); err != nil {
			if errors.Is(err, redis.Nil) {
				utils.SendJSON(w, utils.ApiResponse{
					Error: "url not found",
				}, http.StatusNotFound)
				return
			}

			slog.Error("error delete url", "error", err)
			utils.SendJSON(w, utils.ApiResponse{
				Error: "something went wrong",
			}, http.StatusInternalServerError)
			return
		}

		utils.SendJSON(w, utils.ApiResponse{}, http.StatusNoContent)
	}
}

type updateBody struct {
	NewURL string `json:"new_url"`
}

func HandleUpdateShortenedURL(db repositories.UrlContract) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		var body updateBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.SendJSON(w, utils.ApiResponse{Error: "invalid request body"}, http.StatusUnprocessableEntity)
			return
		}

		if body.NewURL == "" {
			utils.SendJSON(w, utils.ApiResponse{Error: "New URL is required"}, http.StatusBadRequest)
			return
		}

		if _, err := url.Parse(body.NewURL); err != nil {
			utils.SendJSON(w, utils.ApiResponse{Error: "invalid URL"}, http.StatusBadRequest)
			return
		}

		code, err := db.UpdateURL(r.Context(), code, body.NewURL)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				utils.SendJSON(w, utils.ApiResponse{
					Error: "url not found",
				}, http.StatusNotFound)
				return
			}
			slog.Error("error saving url", "error", err)
			utils.SendJSON(w, utils.ApiResponse{
				Error: "something went wrong",
			}, http.StatusInternalServerError)
			return
		}

		utils.SendJSON(w, utils.ApiResponse{Data: code}, http.StatusCreated)
	}
}
