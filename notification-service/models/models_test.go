package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTrendAnalysisEvent_JSONMarshaling(t *testing.T) {
	userID := uuid.New()
	event := TrendAnalysisEvent{
		Event:   "trend_analysis",
		JobType: "daily",
		Status:  "success",
		Date:    "2023-10-01",
		SuccessfulUsers: []SuccessfulUser{
			{UserID: userID, Email: "test@example.com"},
		},
		NotifyAdmin: false,
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal TrendAnalysisEvent: %v", err)
	}

	// Test unmarshaling
	var unmarshaled TrendAnalysisEvent
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal TrendAnalysisEvent: %v", err)
	}

	// Verify fields
	if unmarshaled.Event != event.Event {
		t.Errorf("Expected Event %s, got %s", event.Event, unmarshaled.Event)
	}

	if unmarshaled.JobType != event.JobType {
		t.Errorf("Expected JobType %s, got %s", event.JobType, unmarshaled.JobType)
	}

	if unmarshaled.Status != event.Status {
		t.Errorf("Expected Status %s, got %s", event.Status, unmarshaled.Status)
	}

	if unmarshaled.Date != event.Date {
		t.Errorf("Expected Date %s, got %s", event.Date, unmarshaled.Date)
	}

	if len(unmarshaled.SuccessfulUsers) != 1 {
		t.Errorf("Expected 1 successful user, got %d", len(unmarshaled.SuccessfulUsers))
	}

	if unmarshaled.SuccessfulUsers[0].UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, unmarshaled.SuccessfulUsers[0].UserID)
	}

	if unmarshaled.SuccessfulUsers[0].Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", unmarshaled.SuccessfulUsers[0].Email)
	}
}

func TestTrendAnalysisEvent_FailureEvent(t *testing.T) {
	event := TrendAnalysisEvent{
		Event:         "trend_analysis",
		JobType:       "weekly",
		Status:        "failure",
		Date:          "2023-10-01",
		FailedUserIDs: []string{"user1", "user2", "user3"},
		ErrorSummary:  "Database connection timeout",
		NotifyAdmin:   true,
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal failure event: %v", err)
	}

	// Test unmarshaling
	var unmarshaled TrendAnalysisEvent
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal failure event: %v", err)
	}

	// Verify failure-specific fields
	if len(unmarshaled.FailedUserIDs) != 3 {
		t.Errorf("Expected 3 failed user IDs, got %d", len(unmarshaled.FailedUserIDs))
	}

	if unmarshaled.ErrorSummary != "Database connection timeout" {
		t.Errorf("Expected error summary 'Database connection timeout', got %s", unmarshaled.ErrorSummary)
	}

	if !unmarshaled.NotifyAdmin {
		t.Error("Expected NotifyAdmin to be true")
	}
}

func TestUserNotification_TableName(t *testing.T) {
	notification := UserNotification{}
	tableName := notification.TableName()

	expectedTableName := "user_notifications"
	if tableName != expectedTableName {
		t.Errorf("Expected table name %s, got %s", expectedTableName, tableName)
	}
}

func TestUserNotification_Fields(t *testing.T) {
	userID := uuid.New()
	now := time.Now()
	trendID := uuid.New()

	notification := UserNotification{
		ID:                  uuid.New(),
		UserID:              userID,
		HasNewDailyTrend:    true,
		LastDailyTrendDate:  &now,
		LastDailyTrendID:    &trendID,
		HasNewWeeklyTrend:   false,
		LastWeeklyTrendDate: nil,
		LastWeeklyTrendID:   nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Test field assignments
	if notification.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, notification.UserID)
	}

	if !notification.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be true")
	}

	if notification.HasNewWeeklyTrend {
		t.Error("Expected HasNewWeeklyTrend to be false")
	}

	if notification.LastDailyTrendDate == nil || !notification.LastDailyTrendDate.Equal(now) {
		t.Error("Expected LastDailyTrendDate to be set correctly")
	}

	if notification.LastWeeklyTrendDate != nil {
		t.Error("Expected LastWeeklyTrendDate to be nil")
	}

	if notification.LastDailyTrendID == nil || *notification.LastDailyTrendID != trendID {
		t.Error("Expected LastDailyTrendID to be set correctly")
	}

	if notification.LastWeeklyTrendID != nil {
		t.Error("Expected LastWeeklyTrendID to be nil")
	}
}

func TestUserNotificationResponse_JSONTags(t *testing.T) {
	now := time.Now()
	response := UserNotificationResponse{
		HasNewDailyTrend:    true,
		HasNewWeeklyTrend:   false,
		LastDailyTrendDate:  &now,
		LastWeeklyTrendDate: nil,
	}

	// Test JSON marshaling to verify field names
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal UserNotificationResponse: %v", err)
	}

	// Check that JSON uses camelCase field names
	jsonStr := string(data)

	if !contains(jsonStr, "hasNewDailyTrend") {
		t.Error("Expected JSON to contain hasNewDailyTrend field")
	}

	if !contains(jsonStr, "hasNewWeeklyTrend") {
		t.Error("Expected JSON to contain hasNewWeeklyTrend field")
	}

	if !contains(jsonStr, "lastDailyTrendDate") {
		t.Error("Expected JSON to contain lastDailyTrendDate field")
	}

	if !contains(jsonStr, "lastWeeklyTrendDate") {
		t.Error("Expected JSON to contain lastWeeklyTrendDate field")
	}
}

func TestSuccessfulUser_Fields(t *testing.T) {
	userID := uuid.New()
	user := SuccessfulUser{
		UserID: userID,
		Email:  "test@example.com",
	}

	if user.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, user.UserID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
}

func TestUserFailureInfo_Fields(t *testing.T) {
	userID := uuid.New()
	testError := &testError{"test error"}

	failureInfo := UserFailureInfo{
		UserID: userID,
		Errors: testError,
	}

	if failureInfo.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, failureInfo.UserID)
	}

	if failureInfo.Errors.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got %s", failureInfo.Errors.Error())
	}
}

// Helper functions for testing
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Test error type for UserFailureInfo testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
