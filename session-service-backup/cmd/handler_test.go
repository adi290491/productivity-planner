package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionHandlers(t *testing.T) {
	tests := []struct {
		name               string
		endpoint           string
		method             string
		body               map[string]string
		userID             string
		expectedStatusCode int
	}{
		{
			name:     "StartSession - success",
			endpoint: "/sessions/v1/start-session",
			method:   "POST",
			body: map[string]string{
				"session_type": "focus",
			},
			userID:             "11111111-1111-1111-1111-111111111111",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:     "StartSession - missing user ID",
			endpoint: "/sessions/v1/start-session",
			method:   "POST",
			body: map[string]string{
				"session_type": "focus",
			},
			userID:             "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "StartSession - invalid body",
			endpoint:           "/sessions/v1/start-session",
			method:             "POST",
			body:               nil, // no body sent
			userID:             "11111111-1111-1111-1111-111111111111",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "StartSession - invalid session type",
			endpoint: "/sessions/v1/start-session",
			method:   "POST",
			body: map[string]string{
				"session_type": "invalid",
			},
			userID:             "11111111-1111-1111-1111-111111111111",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "StopSession - success",
			endpoint: "/sessions/v1/stop-session",
			method:   "PATCH",
			body: map[string]string{
				"session_type": "focus",
			},
			userID:             "11111111-1111-1111-1111-111111111111",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:     "StopSession - invalid session type",
			endpoint: "/sessions/v1/stop-session",
			method:   "PATCH",
			body: map[string]string{
				"session_type": "badtype",
			},
			userID:             "11111111-1111-1111-1111-111111111111",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "StopSession - missing user ID",
			endpoint: "/sessions/v1/stop-session",
			method:   "PATCH",
			body: map[string]string{
				"session_type": "focus",
			},
			userID:             "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				jsonBody, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.endpoint, nil)
			}

			if tc.userID != "" {
				req.Header.Set("X-USER-ID", tc.userID)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatusCode {
				t.Errorf("%s: expected status %d, got %d", tc.name, tc.expectedStatusCode, w.Code)
			}
		})
	}
}
