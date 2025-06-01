package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignupHandler(t *testing.T) {

	tests := []struct {
		name               string
		body               map[string]string
		expectedStatusCode int
	}{
		{
			name: "successful signup",
			body: map[string]string{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "missing name",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing email",
			body: map[string]string{
				"name":     "Test User",
				"password": "password123",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: map[string]string{
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		jsonBody, err := json.Marshal(tc.body)
		if err != nil {
			t.Fatalf("Failed to marshal test body: %v", err)
		}

		req, err := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("Expected %d OK, got %d", tc.expectedStatusCode, w.Code)
		}
	}

}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name               string
		body               map[string]string
		expectedStatusCode int
		expectToken        bool
	}{
		{
			name: "successful login",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "1234", // assuming this matches the hash
			},
			expectedStatusCode: http.StatusOK,
			expectToken:        true,
		},
		{
			name: "invalid password",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name: "missing email",
			body: map[string]string{
				"password": "1234",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": "test@example.com",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.body)

			req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatusCode {
				t.Errorf("[%s] expected status %d, got %d", tc.name, tc.expectedStatusCode, w.Code)
			}

			if tc.expectToken {
				var resp map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				if err != nil || resp["token"] == "" {
					t.Errorf("[%s] expected token in response, got: %s", tc.name, w.Body.String())
				}
			}
		})
	}
}
