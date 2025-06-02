package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDailySummary(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		queryDate      string
		expectedStatus int
	}{
		{
			name:           "successful daily summary",
			userID:         "11111111-1111-1111-1111-111111111111",
			queryDate:      "2025-05-25",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user ID",
			userID:         "",
			queryDate:      "2025-05-25",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid date",
			userID:         "11111111-1111-1111-1111-111111111111",
			queryDate:      "invalid-date",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "no sessions found",
			userID:         "notfound",
			queryDate:      "2025-05-25",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/summary/daily?date="+tc.queryDate, nil)
			req.Header.Set("X-USER-ID", tc.userID)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetWeeklySummary(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		startDate      string
		expectedStatus int
	}{
		{
			name:           "successful weekly summary",
			userID:         "11111111-1111-1111-1111-111111111111",
			startDate:      "2025-05-19",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user ID",
			userID:         "",
			startDate:      "2025-05-19",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid start date",
			userID:         "11111111-1111-1111-1111-111111111111",
			startDate:      "invalid-date",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "no sessions found",
			userID:         "notfound",
			startDate:      "2025-05-19",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/summary/weekly?start_date="+tc.startDate, nil)
			req.Header.Set("X-USER-ID", tc.userID)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
