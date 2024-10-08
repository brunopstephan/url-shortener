package api

import (
	"log/slog"
	"net/http"
	"urlShortener/internal/repositories"
	"urlShortener/internal/utils"
)

type getAllUrlsResponse struct {
	URLs map[string]string `json:"urls"`
}

func handleGetAllUrls(db repositories.UrlContract) http.HandlerFunc {
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
