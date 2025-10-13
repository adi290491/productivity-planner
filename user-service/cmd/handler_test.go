package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adi290491/productivity-planner/user-service/user"
	"github.com/adi290491/productivity-planner/user-service/utils"
	"github.com/gin-gonic/gin"
)

// setupTestRouter creates a test router with mock services
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &user.MockUserService{}
	jwtUtil := utils.JWTUtil{Secret: []byte("test-secret")}

	handler := &Handler{
		Svc:     mockService,
		JwtUtil: jwtUtil,
	}

	RegisterEndpoints(router, handler)
	return router
}

func TestSignupHandler(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name               string
		body               map[string]interface{}
		expectedStatusCode int
		expectUserData     bool
	}{
		{
			name: "successful signup",
			body: map[string]interface{}{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatusCode: http.StatusOK,
			expectUserData:     true,
		},
		{
			name: "missing name",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectUserData:     false,
		},
		{
			name: "missing email",
			body: map[string]interface{}{
				"name":     "Test User",
				"password": "password123",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectUserData:     false,
		},
		{
			name: "missing password",
			body: map[string]interface{}{
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectUserData:     false,
		},
		{
			name: "invalid email format",
			body: map[string]interface{}{
				"name":     "Test User",
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectUserData:     false,
		},
		{
			name:               "empty request body",
			body:               map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectUserData:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatusCode, w.Code, w.Body.String())
			}

			if tc.expectUserData && tc.expectedStatusCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				user, exists := response["user"]
				if !exists {
					t.Error("Expected user data in response")
				}

				userMap, ok := user.(map[string]interface{})
				if !ok {
					t.Error("User data should be a map")
				}

				if _, hasID := userMap["id"]; !hasID {
					t.Error("Expected user ID in response")
				}
				if _, hasEmail := userMap["email"]; !hasEmail {
					t.Error("Expected user email in response")
				}
				if _, hasName := userMap["name"]; !hasName {
					t.Error("Expected user name in response")
				}
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name               string
		body               map[string]interface{}
		expectedStatusCode int
		expectToken        bool
	}{
		{
			name: "successful login",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "1234",
			},
			expectedStatusCode: http.StatusOK,
			expectToken:        true,
		},
		{
			name: "invalid password",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name: "nonexistent user",
			body: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "1234",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name: "missing email",
			body: map[string]interface{}{
				"password": "1234",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name: "missing password",
			body: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name: "invalid email format",
			body: map[string]interface{}{
				"email":    "invalid-email",
				"password": "1234",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
		{
			name:               "empty request body",
			body:               map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectToken:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatalf("Failed to marshal test body: %v", err)
			}

			req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatusCode {
				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatusCode, w.Code, w.Body.String())
			}

			if tc.expectToken {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				token, exists := resp["token"]
				if !exists {
					t.Error("Expected token in response")
				}

				tokenStr, ok := token.(string)
				if !ok || tokenStr == "" {
					t.Error("Expected non-empty token string")
				}
			}
		})
	}
}

func TestGetUsersBatchHandler(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name               string
		body               map[string]interface{}
		expectedStatusCode int
		expectUsers        bool
	}{
		{
			name: "successful batch request",
			body: map[string]interface{}{
				"user_ids": []string{
					"11111111-1111-1111-1111-111111111111",
					"22222222-2222-2222-2222-222222222222",
				},
			},
			expectedStatusCode: http.StatusOK,
			expectUsers:        true,
		},
		{
			name: "empty user_ids array",
			body: map[string]interface{}{
				"user_ids": []string{},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectUsers:        false,
		},
		{
			name:               "missing user_ids field",
			body:               map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectUsers:        false,
		},
		{
			name:               "invalid JSON",
			body:               nil, // Will send malformed JSON
			expectedStatusCode: http.StatusBadRequest,
			expectUsers:        false,
		},
		{
			name: "invalid UUID format",
			body: map[string]interface{}{
				"user_ids": []string{"invalid-uuid"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectUsers:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.body == nil {
				// Test malformed JSON
				req, err = http.NewRequest("POST", "/users/batch", strings.NewReader("{invalid json"))
			} else {
				jsonBody, marshalErr := json.Marshal(tc.body)
				if marshalErr != nil {
					t.Fatalf("Failed to marshal test body: %v", marshalErr)
				}
				req, err = http.NewRequest("POST", "/users/batch", bytes.NewBuffer(jsonBody))
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatusCode {
				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatusCode, w.Code, w.Body.String())
			}

			if tc.expectUsers && tc.expectedStatusCode == http.StatusOK {
				var response []map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if len(response) == 0 {
					t.Error("Expected users in response")
				}

				// Verify each user has required fields
				for _, userMap := range response {
					if _, hasID := userMap["id"]; !hasID {
						t.Error("Expected user ID in response")
					}
					if _, hasEmail := userMap["email"]; !hasEmail {
						t.Error("Expected user email in response")
					}
					if _, hasName := userMap["name"]; !hasName {
						t.Error("Expected user name in response")
					}
				}
			}
		})
	}
}

// Test handler error scenarios with service failures
func TestHandlerServiceFailures(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("signup service failure", func(t *testing.T) {
		router := gin.New()
		mockService := &user.MockUserService{ShouldFailSignup: true}
		jwtUtil := utils.JWTUtil{Secret: []byte("test-secret")}
		handler := &Handler{Svc: mockService, JwtUtil: jwtUtil}
		RegisterEndpoints(router, handler)

		body := map[string]interface{}{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "password123",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})

	t.Run("login service failure", func(t *testing.T) {
		router := gin.New()
		mockService := &user.MockUserService{ShouldFailLogin: true}
		jwtUtil := utils.JWTUtil{Secret: []byte("test-secret")}
		handler := &Handler{Svc: mockService, JwtUtil: jwtUtil}
		RegisterEndpoints(router, handler)

		body := map[string]interface{}{
			"email":    "test@example.com",
			"password": "1234",
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("get users batch service failure", func(t *testing.T) {
		router := gin.New()
		mockService := &user.MockUserService{ShouldFailGetUsersBatch: true}
		jwtUtil := utils.JWTUtil{Secret: []byte("test-secret")}
		handler := &Handler{Svc: mockService, JwtUtil: jwtUtil}
		RegisterEndpoints(router, handler)

		body := map[string]interface{}{
			"user_ids": []string{
				"11111111-1111-1111-1111-111111111111",
			},
		}

		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/users/batch", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})
}
