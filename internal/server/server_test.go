package server_test

import (
	"bytes"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"homework/internal/config"
	"homework/internal/middleware"
	"homework/internal/server"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil)
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Test!"))
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	middleware.LoggingMiddleware(handler).ServeHTTP(recorder, req)

	logOutput := buf.String()
	expectedlogOutput := "Incoming Request: method GET, endpoint: /test"
	assert.Contains(t, logOutput, expectedlogOutput)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestBasicAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authorized!"))
	})

	testCases := []struct {
		username     string
		password     string
		expectedCode int
	}{
		{"user", "password", http.StatusOK},
		{"invalid_user", "invalid_password", http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", "/test", nil)
		auth := base64.StdEncoding.EncodeToString([]byte(tc.username + ":" + tc.password))
		req.Header.Set("Authorization", "Basic "+auth)

		recorder := httptest.NewRecorder()
		middleware.BasicAuthMiddleware(handler).ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, tc.expectedCode)

	}
}

func TestCustomMiddlewareAddedToServer(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current directory: %v", err)
	}

	configPath := filepath.Join(currentDir, "..", "..", "config", "local.yaml")

	cfg := config.LoadConfig(configPath)

	srv := server.NewServer(cfg, customMiddleware)

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	err = srv.Shutdown(nil)
	if err != nil {
		t.Errorf("Error shutting down server: %v", err)
	}
}

func customMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		h.ServeHTTP(w, r)
		log.Println(time.Since(now))
	})
}
