package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	DailySummary  *model.DailySessionSummary
	WeeklySummary *model.WeeklySessionSummary
	Err           error
}

func (m *MockService) GetDailySessionSummary(ctx context.Context, userID string, date string) (*model.DailySessionSummary, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if userID == "notfound" {
		return nil, fmt.Errorf("no sessions found for the given day")
	}
	if date == "invalid-date" {
		return nil, fmt.Errorf("invalid date format")
	}
	return m.DailySummary, nil
}

func (m *MockService) GetWeeklySessionSummary(ctx context.Context, userID string, startDate string) (*model.WeeklySessionSummary, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if userID == "notfound" {
		return nil, fmt.Errorf("no sessions found")
	}
	if startDate == "invalid-date" {
		return nil, fmt.Errorf("invalid date format")
	}
	return m.WeeklySummary, nil
}

func TestGetDailySummary_Success(t *testing.T) {
	mockSvc := &MockService{
		DailySummary: &model.DailySessionSummary{
			Date:      "2025-05-25",
			TotalTime: "2h30m",
			Breakdown: map[string]string{
				"focus": "1h30m",
			},
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/daily?date=2025-05-25", nil)
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	handler.GetDailySummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "2025-05-25") {
		t.Errorf("expected date in response, got: %s", body)
	}
}

func TestGetDailySummary_MissingUserID(t *testing.T) {
	mockSvc := &MockService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/daily?date=2025-05-25", nil)
	w := httptest.NewRecorder()

	handler.GetDailySummary(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetDailySummary_InvalidDate(t *testing.T) {
	mockSvc := &MockService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/daily?date=invalid-date", nil)
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	handler.GetDailySummary(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetDailySummary_NoSessions(t *testing.T) {
	mockSvc := &MockService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/daily?date=2025-05-25", nil)
	req.Header.Set("X-USER-ID", "notfound")
	w := httptest.NewRecorder()

	handler.GetDailySummary(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestGetWeeklySummary_Success(t *testing.T) {
	mockSvc := &MockService{
		WeeklySummary: &model.WeeklySessionSummary{
			StartDate: "2025-05-19",
			EndDate:   "2025-05-25",
			TotalTime: "10h0m",
			DailySummaries: []*model.DailySessionSummary{
				{
					Date:      "2025-05-19",
					TotalTime: "2h0m",
					Breakdown: map[string]string{
						"focus": "1h0m",
					},
				},
			},
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/weekly?start_date=2025-05-19", nil)
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	handler.GetWeeklySummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "2025-05-19") {
		t.Errorf("expected start date in response, got: %s", body)
	}
}

func TestGetWeeklySummary_MissingUserID(t *testing.T) {
	mockSvc := &MockService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/weekly?start_date=2025-05-19", nil)
	w := httptest.NewRecorder()

	handler.GetWeeklySummary(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetWeeklySummary_InvalidDate(t *testing.T) {
	mockSvc := &MockService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest("GET", "/summary/weekly?start_date=invalid-date", nil)
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")
	w := httptest.NewRecorder()

	handler.GetWeeklySummary(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	handler := NewHandler(&MockService{})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "healthy") {
		t.Errorf("expected 'healthy' in response, got: %s", body)
	}
}

func TestReady(t *testing.T) {
	handler := NewHandler(&MockService{})

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	handler.Ready(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "ready") {
		t.Errorf("expected 'ready' in response, got: %s", body)
	}
}

