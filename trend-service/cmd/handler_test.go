package main

import (
	"net/http"
	"net/http/httptest"

	models "github.com/adi290491/productivity-planner/trend-service/model"

	"testing"

	"github.com/adi290491/productivity-planner/trend-service/trend"
	"github.com/gin-gonic/gin"
)

func TestHandlers(t *testing.T) {
	var tests = []struct {
		name               string
		url                string
		method             string
		expectedStatusCode int
		userId             string
	}{
		{"GetDailyTrend", "/trend/daily?days=2", "GET", http.StatusOK, "11111111-1111-1111-1111-111111111111"},
		{"GetWeeklyTrend", "/trend/weekly?weeks=3", "GET", http.StatusOK, "11111111-1111-1111-1111-111111111111"},
		{"GetDailyTrendNoQueryParam", "/trend/daily", "GET", http.StatusOK, "11111111-1111-1111-1111-111111111111"},
		{"GetWeeklyTrendNoQueryParam", "/trend/weekly", "GET", http.StatusOK, "11111111-1111-1111-1111-111111111111"},
		{"GetUnviewedTrendsCount", "/trend/unViewed", "GET", http.StatusOK, "11111111-1111-1111-1111-111111111111"},
		{"MarkTrendsAsViewedDaily", "/trend/mark-viewed?type=daily", "POST", http.StatusOK, "1111-1111"},
		{"MarkTrendsAsViewedWeekly", "/trend/mark-viewed?type=weekly", "POST", http.StatusOK, "1111-1111"},
		{"GetDailyTrendUnauthorized", "/trend/daily", "GET", http.StatusUnauthorized, ""},
		{"GetWeeklyTrendUnauthorized", "/trend/weekly", "GET", http.StatusUnauthorized, ""},
		{"GetUnviewedTrendsCountUnauthorized", "/trend/unViewed", "GET", http.StatusUnauthorized, ""},
		{"GetUnviewedTrendsCountInvalidUUID", "/trend/unViewed", "GET", http.StatusInternalServerError, "invalid-uuid"},
		{"MarkTrendsAsViewedUnauthorized", "/trend/mark-viewed?type=daily", "POST", http.StatusUnauthorized, ""},
		{"MarkTrendsAsViewedInvalidType", "/trend/mark-viewed?type=invalid", "POST", http.StatusBadRequest, "11111111-1111-1111-1111-111111111111"},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	RegisterEndpoints(router, &handler)

	// run tests using the above handler
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.url, nil)
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

func TestGetUnviewedTrendsCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	router.GET("/trend/unViewed", handler.GetUnviewedTrendsCount)

	t.Run("returns 200 with valid user id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/unViewed", nil)
		req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is missing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/unViewed", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is empty string", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/unViewed", nil)
		req.Header.Set("X-USER-ID", "")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is whitespace", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/unViewed", nil)
		req.Header.Set("X-USER-ID", "   ")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 500 if user id is invalid UUID format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/trend/unViewed", nil)
		req.Header.Set("X-USER-ID", "invalid-uuid-format")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestMarkTrendsAsViewed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &trend.TrendService{Repo: &models.TestDBRepo{}}
	handler := Handler{svc: mockService}
	router.POST("/trend/mark-viewed", handler.MarkTrendsAsViewed)

	t.Run("returns 200 with valid user id and daily type", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed?type=daily", nil)
		req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 200 with valid user id and weekly type", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed?type=weekly", nil)
		req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is missing", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed?type=daily", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is empty string", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed?type=daily", nil)
		req.Header.Set("X-USER-ID", "")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 401 if user id is whitespace", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed?type=daily", nil)
		req.Header.Set("X-USER-ID", "   ")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid trend type", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed?type=invalid", nil)
		req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 if type parameter is missing", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/trend/mark-viewed", nil)
		req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
