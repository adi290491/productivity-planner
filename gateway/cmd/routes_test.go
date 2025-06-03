package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set a test frontend origin
	os.Setenv("FRONTEND_ORIGIN", "http://localhost:5173")

	router := gin.New()
	RegisterRoutes(router)

	tests := []struct {
		name               string
		method             string
		url                string
		body               string
		headers            map[string]string
		expectedStatusCode int
	}{
		{
			name:               "Signup route works",
			method:             "POST",
			url:                "/users/signup",
			body:               `{"email":"test@example.com","password":"1234","name":"Test"}`,
			expectedStatusCode: http.StatusBadGateway, // no backend, so will 502
		},
		{
			name:               "Login route works",
			method:             "POST",
			url:                "/users/login",
			body:               `{"email":"test@example.com","password":"1234"}`,
			expectedStatusCode: http.StatusBadGateway,
		},
		{
			name:               "Protected route without JWT fails",
			method:             "POST",
			url:                "/sessions/v1/start-session",
			body:               `{"session_type":"focus"}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:   "CORS preflight allowed",
			method: "OPTIONS",
			url:    "/users/signup",
			headers: map[string]string{
				"Origin":                        "http://localhost:5173",
				"Access-Control-Request-Method": "POST",
			},
			expectedStatusCode: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.headers != nil {
				for k, v := range tc.headers {
					req.Header.Set(k, v)
				}
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatusCode {
				t.Errorf("Expected status %d, got %d", tc.expectedStatusCode, w.Code)
			}
		})
	}
}
