package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		name               string
		url                string
		method             string
		body               map[string]string
		expectedStatusCode int
	}{
		{
			name:   "CreateUser",
			url:    "/users/signup",
			method: "POST",
			body: map[string]string{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "Test@123",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "GetUser",
			url:    "/users/login",
			method: "POST",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "Test@123",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.body)
			req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d for %s, got %d", test.expectedStatusCode, test.url, w.Code)
			}
		})
	}
}
