package main

import (
	"net/http"
	"net/http/httptest"
	models "productivity-planner/trend-service/model"
	"productivity-planner/trend-service/trend"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandlers(t *testing.T) {
	var tests = []struct {
		name               string
		url                string
		expectedStatusCode int
		userId             string
	}{
		{"GetDailyTrend", "/trend/daily?days=2", http.StatusOK, "1111-1111"},
		{"GetWeeklyTrend", "/trend/weekly?weeks=3", http.StatusOK, "1111-1111"},
		{"GetDailyTrendNoQueryParam", "/trend/daily", http.StatusOK, "1111-1111"},
		{"GetWeeklyTrendNoQueryParam", "/trend/weekly", http.StatusOK, "1111-1111"},
		{"GetDailyTrendUnauthorized", "/trend/daily", http.StatusUnauthorized, ""},
		{"GetWeeklyTrendUnauthorized", "/trend/weekly", http.StatusUnauthorized, ""},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	RegisterEndpoints(router, &handler)

	// run tests using the above handler
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			req.Header.Set("X-USER-ID", test.userId)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d for %s, got %d", test.expectedStatusCode, test.url, w.Code)
			}
		})
	}
}
func TestGetDailyTrend(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock TrendService
	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	router.GET("/trend/daily", handler.GetDailyTrend)

	t.Run("returns 200 with valid user id and days param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/daily?days=2", nil)
		req.Header.Set("X-USER-ID", "test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 200 with valid user id and no days param (default)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/daily", nil)
		req.Header.Set("X-USER-ID", "test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is missing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/daily?days=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
func TestGetDailyTrend_UserIdMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	router.GET("/trend/daily", handler.GetDailyTrend)

	req, _ := http.NewRequest("GET", "/trend/daily?days=2", nil)
	// Do not set X-USER-ID header to trigger the if userId == "" branch
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
func TestGetWeeklyTrend(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock TrendService
	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	router.GET("/trend/weekly", handler.GetWeeklyTrend)

	t.Run("returns 200 with valid user id and weeks param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/weekly?weeks=2", nil)
		req.Header.Set("X-USER-ID", "test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 200 with valid user id and no weeks param (default)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/weekly", nil)
		req.Header.Set("X-USER-ID", "test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is missing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/weekly?weeks=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}

func TestGetWeeklyTrend_UserIdMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	router.GET("/trend/weekly", handler.GetWeeklyTrend)

	req, _ := http.NewRequest("GET", "/trend/weekly?weeks=2", nil)
	// Do not set X-USER-ID header to trigger the if userId == "" branch
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
