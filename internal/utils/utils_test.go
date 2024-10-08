package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenCode(t *testing.T) {
	code := GenCode()

	if code == "" {
		t.Error("expected code to be non-empty")
	}

	if len(code) != 8 {
		t.Errorf("expected code to have length 8, got %d", len(code))
	}
}

// reason why it doesn't have a fail in marshal test case is because the
// function expect a specific struct to be passed as a parameter
// so it can't break the marshal
func TestSendJSON(t *testing.T) {
	tests := []struct {
		name           string
		resp           ApiResponse
		status         int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid response",
			resp:           ApiResponse{Data: "success"},
			status:         http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SendJSON(w, tt.resp, tt.status)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body := new(bytes.Buffer)
			body.ReadFrom(res.Body)
			if body.String() != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body.String())
			}
		})
	}
}
