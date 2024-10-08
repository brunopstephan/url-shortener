package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlShortener/internal/utils"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUrlRepository struct {
	mock.Mock
}

func (m *MockUrlRepository) SaveShortenedURL(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (m *MockUrlRepository) GetURL(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

func (m *MockUrlRepository) GetAllURL(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestPostShortenedURL_ValidRequest(t *testing.T) {
	validUrl := "https://example.com"
	tt := struct {
		body           postBody
		mockSaveReturn string
		mockSaveError  error
		expectedCode   int
		expectedBody   utils.ApiResponse
	}{
		body:           postBody{URL: validUrl},
		mockSaveReturn: validUrl,
		mockSaveError:  nil,
		expectedCode:   http.StatusCreated,
		expectedBody:   utils.ApiResponse{Data: validUrl},
	}
	mockStore := new(MockUrlRepository)
	mockStore.On("SaveShortenedURL", mock.Anything, tt.body.URL).Return(tt.mockSaveReturn, tt.mockSaveError)
	handler := HandlePostShortenedURL(mockStore)

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(tt.body)

	req := httptest.NewRequest("POST", "/api/shorten", &requestBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	var actualResponse utils.ApiResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, tt.expectedBody, actualResponse)

	mockStore.AssertExpectations(t)
}

func TestPostShortenedURL_MissingParams(t *testing.T) {
	tt := struct {
		name         string
		body         postBody
		expectedCode int
		expectedBody utils.ApiResponse
	}{
		body:         postBody{},
		expectedCode: http.StatusBadRequest,
		expectedBody: utils.ApiResponse{Error: "URL is required"},
	}

	mockStore := new(MockUrlRepository)
	handler := HandlePostShortenedURL(mockStore)

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(tt.body)

	req := httptest.NewRequest("POST", "/api/shorten", &requestBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	var actualResponse utils.ApiResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, tt.expectedBody, actualResponse)
}

func TestPostShortenedURL_InvalidRequest(t *testing.T) {
	tt := struct {
		body         string
		expectedCode int
		expectedBody utils.ApiResponse
	}{
		body:         "",
		expectedCode: http.StatusUnprocessableEntity,
		expectedBody: utils.ApiResponse{Error: "invalid request body"},
	}

	mockStore := new(MockUrlRepository)
	handler := HandlePostShortenedURL(mockStore)

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(tt.body)

	req := httptest.NewRequest("POST", "/api/shorten", &requestBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	var actualResponse utils.ApiResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, tt.expectedBody, actualResponse)
}

func TestPostShortenedURL_DatabaseError(t *testing.T) {
	tt := struct {
		body         postBody
		expectedCode int
		expectedBody utils.ApiResponse
	}{
		body:         postBody{URL: "https://example.com"},
		expectedCode: http.StatusInternalServerError,
		expectedBody: utils.ApiResponse{
			Error: "something went wrong",
		},
	}

	mockStore := new(MockUrlRepository)
	mockStore.On("SaveShortenedURL", mock.Anything, mock.Anything).Return("", assert.AnError)
	handler := HandlePostShortenedURL(mockStore)

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(tt.body)

	req := httptest.NewRequest("POST", "/api/shorten", &requestBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	var actualResponse utils.ApiResponse
	json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Equal(t, tt.expectedBody, actualResponse)
}

func TestGetShortenedURL_ValidRequest(t *testing.T) {
	validUrl := "https://example.com"
	tt := struct {
		mockSaveReturn string
		mockSaveError  error
		expectedCode   int
		expectedBody   utils.ApiResponse
	}{
		mockSaveReturn: validUrl,
		mockSaveError:  nil,
		expectedCode:   http.StatusOK,
		expectedBody: utils.ApiResponse{
			Data: getShortenedURLResponse{URL: validUrl},
		},
	}
	mockStore := new(MockUrlRepository)
	mockStore.On("GetURL", context.Background(), "").Return(tt.mockSaveReturn, tt.mockSaveError)
	handler := HandleGetShortenedURL(mockStore)

	req := httptest.NewRequest("GET", "/api/123?json=true", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	expectedBody, _ := json.Marshal(tt.expectedBody)

	assert.JSONEq(t, string(expectedBody), w.Body.String())

	mockStore.AssertExpectations(t)
}

func TestGetShortenedURL_UrlNotFound(t *testing.T) {
	tt := struct {
		expectedCode int
		expectedBody utils.ApiResponse
	}{
		expectedCode: http.StatusNotFound,
		expectedBody: utils.ApiResponse{
			Error: "url not found",
		},
	}
	mockStore := new(MockUrlRepository)
	mockStore.On("GetURL", context.Background(), "").Return("", redis.Nil)
	handler := HandleGetShortenedURL(mockStore)

	req := httptest.NewRequest("GET", "/api/123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	expectedBody, _ := json.Marshal(tt.expectedBody)

	assert.JSONEq(t, string(expectedBody), w.Body.String())

	mockStore.AssertExpectations(t)
}

func TestGetShortenedURL_SomethingWentWrong(t *testing.T) {
	tt := struct {
		expectedCode int
		expectedBody utils.ApiResponse
	}{
		expectedCode: http.StatusInternalServerError,
		expectedBody: utils.ApiResponse{
			Error: "something went wrong",
		},
	}
	mockStore := new(MockUrlRepository)
	mockStore.On("GetURL", context.Background(), "").Return("", assert.AnError)
	handler := HandleGetShortenedURL(mockStore)

	req := httptest.NewRequest("GET", "/api/123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	expectedBody, _ := json.Marshal(tt.expectedBody)

	assert.JSONEq(t, string(expectedBody), w.Body.String())

	mockStore.AssertExpectations(t)
}

func TestGetAllURL_ValidRequest(t *testing.T) {
	tt := struct {
		mockSaveReturn map[string]string
		mockSaveError  error
		expectedCode   int
		expectedBody   utils.ApiResponse
	}{
		mockSaveReturn: map[string]string{"123": "https://example.com"},
		mockSaveError:  nil,
		expectedCode:   http.StatusOK,
		expectedBody: utils.ApiResponse{
			Data: getAllUrlsResponse{URLs: map[string]string{"123": "https://example.com"}},
		},
	}
	mockStore := new(MockUrlRepository)
	mockStore.On("GetAllURL", context.Background()).Return(tt.mockSaveReturn, tt.mockSaveError)
	handler := HandleGetAllUrls(mockStore)

	req := httptest.NewRequest("GET", "/dashboard/all", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	expectedBody, _ := json.Marshal(tt.expectedBody)

	assert.JSONEq(t, string(expectedBody), w.Body.String())

	mockStore.AssertExpectations(t)
}

func TestGetAllURL_SomethingWentWrong(t *testing.T) {
	tt := struct {
		expectedCode int
		expectedBody utils.ApiResponse
	}{
		expectedCode: http.StatusInternalServerError,
		expectedBody: utils.ApiResponse{
			Error: "something went wrong",
		},
	}
	mockStore := new(MockUrlRepository)
	mockStore.On("GetAllURL", context.Background()).Return(map[string]string{}, assert.AnError)
	handler := HandleGetAllUrls(mockStore)

	req := httptest.NewRequest("GET", "/dashboard/all", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, tt.expectedCode, w.Code)

	expectedBody, _ := json.Marshal(tt.expectedBody)

	assert.JSONEq(t, string(expectedBody), w.Body.String())

	mockStore.AssertExpectations(t)
}
