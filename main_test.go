package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdateHandler_Auth(t *testing.T) {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}
	logHandler := slog.NewJSONHandler(os.Stdout, logOpts)
	slog.SetDefault(slog.New(logHandler))

	config := &Config{
		Commands: []string{"echo 'test'"},
	}

	os.Setenv("RUNNER_AUTH_TOKEN", "testtoken")
	defer os.Unsetenv("RUNNER_AUTH_TOKEN")

	handler := updateHandler(config)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"Valid Token", "Bearer testtoken", http.StatusOK},
		{"Invalid Token", "Bearer wrongtoken", http.StatusUnauthorized},
		{"No Token", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/update", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
		})
	}
}
