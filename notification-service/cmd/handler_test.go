package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
	"github.com/adi290491/productivity-planner/notification-service/notification"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Mock NotificationService for testing
type MockNotificationService struct {
	processDailyTrendError      error
	processWeeklyTrendError     error
	getUserNotificationResponse *models.UserNotificationResponse
	getUserNotificationError    error
	markNotificationAsReadError error
	stats                       MockProcessingStats
}

func (m *MockNotificationService) ProcessDailyTrendNotifications(ctx context.Context) error {
	return m.processDailyTrendError
}

func (m *MockNotificationService) ProcessWeeklyTrendNotifications(ctx context.Context) error {
	return m.processWeeklyTrendError
}

func (m *MockNotificationService) GetUserNotification(userID uuid.UUID) (*models.UserNotificationResponse, error) {
	return m.getUserNotificationResponse, m.getUserNotificationError
}

func (m *MockNotificationService) MarkNotificationAsRead(userID uuid.UUID, notificationType string) error {
	return m.markNotificationAsReadError
}

func (m *MockNotificationService) GetStats() notification.ProcessingStats {
	return notification.ProcessingStats{
		MessagesProcessed: m.stats.MessagesProcessed,
		LastProcessed:     m.stats.LastProcessed,
	}
}

type MockProcessingStats struct {
	MessagesProcessed int       `json:"messagesProcessed"`
	LastProcessed     time.Time `json:"lastProcessed"`
}

func setupTestHandler(mockService *MockNotificationService) (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	config := &config.AppConfig{
		Profile: "test",
	}

	// Create handler with interface that matches our mock
	handler := &Handler{
		config:              config,
		notificationService: mockService,
	}

	return handler, router
}

func TestProcessDailyTrends_Success(t *testing.T) {
	mockService := &MockNotificationService{
		stats: MockProcessingStats{
			MessagesProcessed: 5,
			LastProcessed:     time.Now(),
		},
	}
	handler, router := setupTestHandler(mockService)

	router.POST("/process/daily", handler.processDailyTrends)

	req, _ := http.NewRequest("POST", "/process/daily", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Error("Expected success to be true")
	}

	if response["message"] != "Daily trend notifications processed successfully" {
		t.Error("Expected success message")
	}
}

func TestProcessDailyTrends_Error(t *testing.T) {
	mockService := &MockNotificationService{
		processDailyTrendError: fmt.Errorf("processing failed"),
	}
	handler, router := setupTestHandler(mockService)

	router.POST("/process/daily", handler.processDailyTrends)

	req, _ := http.NewRequest("POST", "/process/daily", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != false {
		t.Error("Expected success to be false")
	}

	if response["error"] != "processing failed" {
		t.Error("Expected error message")
	}
}

func TestProcessWeeklyTrends_Success(t *testing.T) {
	mockService := &MockNotificationService{
		stats: MockProcessingStats{
			MessagesProcessed: 3,
			LastProcessed:     time.Now(),
		},
	}
	handler, router := setupTestHandler(mockService)

	router.POST("/process/weekly", handler.processWeeklyTrends)

	req, _ := http.NewRequest("POST", "/process/weekly", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProcessWeeklyTrends_Error(t *testing.T) {
	mockService := &MockNotificationService{
		processWeeklyTrendError: fmt.Errorf("weekly processing failed"),
	}
	handler, router := setupTestHandler(mockService)

	router.POST("/process/weekly", handler.processWeeklyTrends)

	req, _ := http.NewRequest("POST", "/process/weekly", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetUserNotifications_Success(t *testing.T) {
	userID := uuid.New()
	mockResponse := &models.UserNotificationResponse{
		HasNewDailyTrend:  true,
		HasNewWeeklyTrend: false,
	}

	mockService := &MockNotificationService{
		getUserNotificationResponse: mockResponse,
	}
	handler, router := setupTestHandler(mockService)

	router.GET("/notifications", handler.getUserNotifications)

	// Set up authenticated user context
	req, _ := http.NewRequest("GET", fmt.Sprintf("/notifications?user_id=%s", userID.String()), nil)
	w := httptest.NewRecorder()

	// Create context with authenticated user
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.getUserNotifications(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.UserNotificationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be true")
	}
}

func TestGetUserNotifications_NoAuth(t *testing.T) {
	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("GET", "/notifications?user_id=123", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// Don't set userID in context - simulate no auth

	handler.getUserNotifications(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetUserNotifications_WrongUser(t *testing.T) {
	authenticatedUserID := uuid.New()
	requestedUserID := uuid.New()

	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/notifications?user_id=%s", requestedUserID.String()), nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", authenticatedUserID)

	handler.getUserNotifications(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetUserNotifications_InvalidUserID(t *testing.T) {
	userID := uuid.New()

	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("GET", "/notifications?user_id=invalid-uuid", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.getUserNotifications(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMarkNotificationAsRead_Success(t *testing.T) {
	userID := uuid.New()

	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/notifications/read?user_id=%s&type=daily", userID.String()), nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.markNotificationAsRead(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Error("Expected success to be true")
	}
}

func TestMarkNotificationAsRead_InvalidType(t *testing.T) {
	userID := uuid.New()

	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/notifications/read?user_id=%s&type=invalid", userID.String()), nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.markNotificationAsRead(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMarkNotificationAsRead_CaseInsensitive(t *testing.T) {
	userID := uuid.New()

	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	// Test with uppercase
	req, _ := http.NewRequest("POST", fmt.Sprintf("/notifications/read?user_id=%s&type=DAILY", userID.String()), nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.markNotificationAsRead(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMarkNotificationAsRead_MissingType(t *testing.T) {
	userID := uuid.New()

	mockService := &MockNotificationService{}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/notifications/read?user_id=%s", userID.String()), nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.markNotificationAsRead(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMarkNotificationAsRead_ServiceError(t *testing.T) {
	userID := uuid.New()

	mockService := &MockNotificationService{
		markNotificationAsReadError: fmt.Errorf("database error"),
	}
	handler, _ := setupTestHandler(mockService)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/notifications/read?user_id=%s&type=daily", userID.String()), nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID)

	handler.markNotificationAsRead(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
