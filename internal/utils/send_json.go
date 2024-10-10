package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ApiResponse struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func SendJSON(w http.ResponseWriter, resp ApiResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("error marshaling response", "error", err)
		SendJSON(w, ApiResponse{Error: "something went wrong sending json (marshal)"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if status == http.StatusNoContent || (resp.Data == nil && resp.Error == "") {
		return
	}

	if _, err := w.Write(data); err != nil {
		SendJSON(w, ApiResponse{Error: "something went wrong sending json (wrinting response)"}, http.StatusInternalServerError)
		return
	}

}
