package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/api/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMaxBodySize(t *testing.T) {
	// Тестовый хендлер
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	// Применяем middleware с лимитом 100 байт
	handler := middleware.MaxBodySize(100)(next)

	tests := []struct {
		name       string
		bodySize   int
		wantStatus int
	}{
		{
			name:       "body within limit",
			bodySize:   50,
			wantStatus: http.StatusOK,
		},
		{
			name:       "body exactly at limit",
			bodySize:   100,
			wantStatus: http.StatusOK,
		},
		{
			name:       "body exceeds limit",
			bodySize:   200,
			wantStatus: http.StatusBadRequest, // MaxBytesReader возвращает 400
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewReader(make([]byte, tt.bodySize))
			req := httptest.NewRequest("POST", "/test", body)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestMaxBodySize_WithCustomErrorHandler(t *testing.T) {
	// Тестовый хендлер с кастомной обработкой ошибки
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			// Если ошибка - возвращаем 413
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.MaxBodySize(100)(next)

	body := bytes.NewReader(make([]byte, 200))
	req := httptest.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestRateLimit(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RateLimit(1, 2)(next) // 1 req/sec, burst 2

	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
