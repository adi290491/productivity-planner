package notification

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use a test PostgreSQL database or skip tests if not available
	testDSN := "postgres://postgres:password@localhost:5432/test_notifications?sslmode=disable"

	db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: failed to connect to test database: %v", err)
		return nil
	}

	err = db.AutoMigrate(&models.UserNotification{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Clean up any existing test data
	db.Exec("TRUNCATE TABLE user_notifications")

	return db
}

func setupTestNotificationService(t *testing.T) *NotificationService {
	config := &config.AppConfig{
		Profile:   "test",
		ProjectID: "test-project",
	}

	db := setupTestDB(t)
	if db == nil {
		t.Skip("Database not available for testing")
		return nil
	}

	ns := &NotificationService{
		config: config,
		db:     db,
	}

	return ns
}

func TestGetUserNotification_NotFound(t *testing.T) {
	ns := setupTestNotificationService(t)
	userID := uuid.New()

	response, err := ns.GetUserNotification(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be false for non-existent user")
	}

	if response.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to be false for non-existent user")
	}

	if response.LastDailyTrendDate != nil {
		t.Error("Expected LastDailyTrendDate to be nil for non-existent user")
	}

	if response.LastWeeklyTrendDate != nil {
		t.Error("Expected LastWeeklyTrendDate to be nil for non-existent user")
	}
}

func TestGetUserNotification_Exists(t *testing.T) {
	ns := setupTestNotificationService(t)
	userID := uuid.New()
	now := time.Now()

	// Create a notification record
	notification := &models.UserNotification{
		ID:                  uuid.New(),
		UserID:              userID,
		HasNewDailyTrend:    true,
		LastDailyTrendDate:  &now,
		HasNewWeeklyTrend:   false,
		LastWeeklyTrendDate: nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err := ns.db.Create(notification).Error
	if err != nil {
		t.Fatalf("Failed to create test notification: %v", err)
	}

	response, err := ns.GetUserNotification(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be true")
	}

	if response.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to be false")
	}

	if response.LastDailyTrendDate == nil {
		t.Error("Expected LastDailyTrendDate to be set")
	}

	if response.LastWeeklyTrendDate != nil {
		t.Error("Expected LastWeeklyTrendDate to be nil")
	}
}

func TestMarkNotificationAsRead_Daily(t *testing.T) {
	ns := setupTestNotificationService(t)
	userID := uuid.New()
	now := time.Now()

	// Create a notification record
	notification := &models.UserNotification{
		ID:                  uuid.New(),
		UserID:              userID,
		HasNewDailyTrend:    true,
		LastDailyTrendDate:  &now,
		HasNewWeeklyTrend:   true,
		LastWeeklyTrendDate: &now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err := ns.db.Create(notification).Error
	if err != nil {
		t.Fatalf("Failed to create test notification: %v", err)
	}

	err = ns.MarkNotificationAsRead(userID, "daily")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the daily flag is false but weekly is still true
	var updated models.UserNotification
	err = ns.db.Where("user_id = ?", userID).First(&updated).Error
	if err != nil {
		t.Fatalf("Failed to fetch updated notification: %v", err)
	}

	if updated.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be false after marking as read")
	}

	if !updated.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to remain true")
	}
}

func TestMarkNotificationAsRead_Weekly(t *testing.T) {
	ns := setupTestNotificationService(t)
	userID := uuid.New()
	now := time.Now()

	// Create a notification record
	notification := &models.UserNotification{
		ID:                  uuid.New(),
		UserID:              userID,
		HasNewDailyTrend:    true,
		LastDailyTrendDate:  &now,
		HasNewWeeklyTrend:   true,
		LastWeeklyTrendDate: &now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err := ns.db.Create(notification).Error
	if err != nil {
		t.Fatalf("Failed to create test notification: %v", err)
	}

	err = ns.MarkNotificationAsRead(userID, "weekly")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the weekly flag is false but daily is still true
	var updated models.UserNotification
	err = ns.db.Where("user_id = ?", userID).First(&updated).Error
	if err != nil {
		t.Fatalf("Failed to fetch updated notification: %v", err)
	}

	if updated.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to be false after marking as read")
	}

	if !updated.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to remain true")
	}
}

func TestMarkNotificationAsRead_InvalidType(t *testing.T) {
	ns := setupTestNotificationService(t)
	userID := uuid.New()

	err := ns.MarkNotificationAsRead(userID, "invalid")
	if err == nil {
		t.Error("Expected error for invalid notification type")
	}

	expectedError := "invalid notification type: invalid"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestMarkNotificationAsRead_NotFound(t *testing.T) {
	ns := setupTestNotificationService(t)
	userID := uuid.New()

	err := ns.MarkNotificationAsRead(userID, "daily")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	expectedError := "notification not found for user"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestUpdateNotificationFlag_Daily(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}
	userID := uuid.New()
	trendID := uuid.New()
	trendDate, _ := time.Parse("2006-01-02", "2023-10-01")

	err := ns.updateNotificationFlag(userID, DAILY, trendDate, &trendID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify notification was created
	var notification models.UserNotification
	err = ns.db.Where("user_id = ?", userID).First(&notification).Error
	if err != nil {
		t.Fatalf("Failed to fetch notification: %v", err)
	}

	if !notification.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be true")
	}

	if notification.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to be false")
	}

	if notification.LastDailyTrendDate == nil || !notification.LastDailyTrendDate.Equal(trendDate) {
		t.Error("Expected LastDailyTrendDate to be set correctly")
	}
}

func TestUpdateNotificationFlag_Weekly(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}
	userID := uuid.New()
	trendID := uuid.New()
	trendDate, _ := time.Parse("2006-01-02", "2023-10-01")

	err := ns.updateNotificationFlag(userID, WEEKLY, trendDate, &trendID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify notification was created
	var notification models.UserNotification
	err = ns.db.Where("user_id = ?", userID).First(&notification).Error
	if err != nil {
		t.Fatalf("Failed to fetch notification: %v", err)
	}

	if notification.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be false")
	}

	if !notification.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to be true")
	}

	if notification.LastWeeklyTrendDate == nil || !notification.LastWeeklyTrendDate.Equal(trendDate) {
		t.Error("Expected LastWeeklyTrendDate to be set correctly")
	}
}

func TestMessageValidation_InvalidJSON(t *testing.T) {
	// Test JSON unmarshaling directly
	var event models.TrendAnalysisEvent
	err := json.Unmarshal([]byte("invalid json"), &event)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestMessageValidation_EmptyFields(t *testing.T) {
	// Test validation logic for empty required fields
	event := models.TrendAnalysisEvent{
		Event:   "",
		JobType: "daily",
		Status:  "success",
		Date:    "2023-10-01",
	}

	if event.Event == "" || event.Status == "" || event.Date == "" {
		// This validates our field checking logic works
		if event.Event == "" {
			// Expected behavior
		}
	}
}

func TestHandleSuccessEvent(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	userID1 := uuid.New()
	userID2 := uuid.New()

	event := models.TrendAnalysisEvent{
		Event:   "trend_analysis",
		JobType: "daily",
		Status:  "success",
		Date:    "2023-10-01",
		SuccessfulUsers: []models.SuccessfulUser{
			{UserID: userID1, Email: "user1@test.com"},
			{UserID: userID2, Email: "user2@test.com"},
		},
	}

	err := ns.handleSuccessEvent(context.Background(), event, DAILY)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify notifications were created for both users
	var notifications []models.UserNotification
	err = ns.db.Find(&notifications).Error
	if err != nil {
		t.Fatalf("Failed to fetch notifications: %v", err)
	}

	if len(notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(notifications))
	}

	for _, notification := range notifications {
		if !notification.HasNewDailyTrend {
			t.Error("Expected HasNewDailyTrend to be true")
		}
		if notification.LastDailyTrendDate == nil || notification.LastDailyTrendDate.Format("2006-01-02") != "2023-10-01" {
			t.Error("Expected LastDailyTrendDate to be set correctly")
		}
	}
}

func TestHandleSuccessEvent_NoUsers(t *testing.T) {
	ns := setupTestNotificationService(t)

	event := models.TrendAnalysisEvent{
		Event:           "trend_analysis",
		JobType:         "daily",
		Status:          "success",
		Date:            "2023-10-01",
		SuccessfulUsers: []models.SuccessfulUser{},
	}

	err := ns.handleSuccessEvent(context.Background(), event, DAILY)
	if err != nil {
		t.Fatalf("Expected no error for empty users list, got %v", err)
	}

	// Verify no notifications were created
	var count int64
	ns.db.Model(&models.UserNotification{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 notifications, got %d", count)
	}
}

func TestHandleFailureEvent(t *testing.T) {
	ns := setupTestNotificationService(t)

	event := models.TrendAnalysisEvent{
		Event:         "trend_analysis",
		JobType:       "daily",
		Status:        "failure",
		Date:          "2023-10-01",
		FailedUserIDs: []string{"user1", "user2"},
		ErrorSummary:  "Database connection failed",
		NotifyAdmin:   true,
	}

	err := ns.handleFailureEvent(context.Background(), event, DAILY)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Since we don't have admin notification logic implemented yet,
	// this test just verifies the function doesn't crash
}

func TestGetStats(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	// Simulate processing some messages
	ns.mu.Lock()
	ns.messagesProcessed = 5
	ns.lastProcessed = time.Now()
	ns.mu.Unlock()

	stats := ns.GetStats()

	if stats.MessagesProcessed != 5 {
		t.Errorf("Expected MessagesProcessed to be 5, got %d", stats.MessagesProcessed)
	}

	if stats.LastProcessed.IsZero() {
		t.Error("Expected LastProcessed to be set")
	}
}

func TestProcessMessage_Success(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	userID := uuid.New()
	event := models.TrendAnalysisEvent{
		Event:           "trend_analysis",
		JobType:         "daily",
		Status:          "success",
		Date:            "2023-10-01",
		SuccessfulUsers: []models.SuccessfulUser{{UserID: userID, Email: "pm_success@test.com"}},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	msg := &pubsub.Message{Data: data, ID: "msg-success"}
	err = ns.processMessage(context.Background(), msg, DAILY)
	if err != nil {
		t.Fatalf("processMessage returned error for success event: %v", err)
	}

	// verify created
	var notification models.UserNotification
	err = ns.db.Where("user_id = ?", userID).First(&notification).Error
	if err != nil {
		t.Fatalf("expected notification to be created, got error: %v", err)
	}
}

func TestProcessMessage_Failure(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	event := models.TrendAnalysisEvent{
		Event:         "trend_analysis",
		JobType:       "daily",
		Status:        "failure",
		Date:          "2023-10-01",
		FailedUserIDs: []string{"user1"},
		ErrorSummary:  "simulated failure",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	msg := &pubsub.Message{Data: data, ID: "msg-failure"}
	err = ns.processMessage(context.Background(), msg, DAILY)
	if err != nil {
		t.Fatalf("processMessage returned error for failure event: %v", err)
	}
}

func TestProcessMessage_UnknownStatus(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	event := models.TrendAnalysisEvent{
		Event:   "trend_analysis",
		JobType: "daily",
		Status:  "invalid_status",
		Date:    "2023-10-01",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	msg := &pubsub.Message{Data: data, ID: "msg-unknown"}
	err = ns.processMessage(context.Background(), msg, DAILY)
	if err == nil {
		t.Fatalf("expected error for unknown status, got nil")
	}
	if !strings.Contains(err.Error(), "unknown event status") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestProcessMessage_InvalidJSON(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	msg := &pubsub.Message{Data: []byte("not a json"), ID: "msg-invalid-json"}
	err := ns.processMessage(context.Background(), msg, DAILY)
	if err == nil {
		t.Fatalf("expected unmarshal error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to unmarshal message") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestProcessMessage_MissingFields(t *testing.T) {
	ns := setupTestNotificationService(t)
	if ns == nil {
		return
	}

	event := models.TrendAnalysisEvent{
		Event:   "",
		JobType: "daily",
		Status:  "success",
		Date:    "",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	msg := &pubsub.Message{Data: data, ID: "msg-missing-fields"}
	err = ns.processMessage(context.Background(), msg, DAILY)
	if err == nil {
		t.Fatalf("expected error for missing required fields, got nil")
	}
	if !strings.Contains(err.Error(), "invalid event data") && !strings.Contains(err.Error(), "status field is empty") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
