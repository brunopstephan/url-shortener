package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"url-shortener/internal/repositories"
	"url-shortener/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type getShortenedURLResponse struct {
	URL string `json:"url"`
}

// HandleGetShortenedURL godoc
// @Summary Get shortened URL
// @Description Get the original URL from the shortened code
// @Tags API
// @Param code path string true "Shortened URL code"
// @Param json query string false "Return JSON response"
// @Success 200 {object} utils.ApiResponse{data=getShortenedURLResponse}
// @Failure 404 {object} utils.ApiResponse{error=string}
// @Failure 500 {object} utils.ApiResponse{error=string}
// @Router /api/{code} [get]
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

// HandlePostShortenedURL godoc
// @Summary Post shortened URL
// @Description Get the original URL from the shortened code
// @Tags API
// @Param data body postBody true "Shortened URL Post Body"
// @Success 201 {object} utils.ApiResponse{data=string}
// @Failure 400 {object} utils.ApiResponse{error=string}
// @Failure 500 {object} utils.ApiResponse{error=string}
// @Failure 422 {object} utils.ApiResponse{error=string}
// @Router /api/shorten [post]
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

// HandleGetAllUrls godoc
// @Summary Get all shortened URL
// @Description Get all shortened URLs and respective codes
// @Security BasicAuth
// @Tags ADMIN
// @Param Authorization header string true "Basic Auth"
// @Success 200 {object} utils.ApiResponse{data=getAllUrlsResponse}
// @Failure 500 {object} utils.ApiResponse{error=string}
// @Failure 401
// @Router /admin/all [get]
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

// HandleDeleteShortenedURL godoc
// @Summary Delete shortened URL
// @Description Delete shortened URL that match the code passed
// @Security BasicAuth
// @Tags ADMIN
// @Param Authorization header string true "Basic Auth"
// @Param code path string true "Shortened URL code"
// @Success 204 {object} utils.ApiResponse{}
// @Failure 500 {object} utils.ApiResponse{error=string}
// @Failure 404 {object} utils.ApiResponse{error=string}
// @Failure 401
// @Router /admin/{code} [delete]
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

// HandleUpdateShortenedURL godoc
// @Summary Update shortened URL
// @Description Update shortened URL that match the code passed
// @Security BasicAuth
// @Tags ADMIN
// @Param Authorization header string true "Basic Auth"
// @Param code path string true "Shortened URL code"
// @Param data body updateBody true "Shortened URL Update Body"
// @Success 201 {object} utils.ApiResponse{data=string}
// @Failure 500 {object} utils.ApiResponse{error=string}
// @Failure 404 {object} utils.ApiResponse{error=string}
// @Failure 422 {object} utils.ApiResponse{error=string}
// @Failure 400 {object} utils.ApiResponse{error=string}
// @Failure 401
// @Router /admin/{code} [put]
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
