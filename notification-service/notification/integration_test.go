package notification

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/pubsub/v2"
	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	testProjectID = "test-project"
	testSubID     = "test-daily-subscription"
)

// Integration test for PubSub subscription functionality
func TestPubSubSubscription_Integration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test: set RUN_INTEGRATION_TESTS=true to run")
	}

	// Set up test database
	testDB := setupIntegrationTestDB(t)
	if testDB == nil {
		t.Skip("Database not available for integration testing")
		return
	}

	// Set up test configuration
	testConfig := &config.AppConfig{
		Profile:            "test",
		ProjectID:          testProjectID,
		DailySubscription:  testSubID,
		WeeklySubscription: "test-weekly-subscription",
	}

	ctx := context.Background()

	// Set up PubSub client (assumes emulator is running)
	pubsubClient, err := pubsub.NewClient(ctx, testProjectID, option.WithoutAuthentication())
	if err != nil {
		t.Skipf("PubSub emulator not available: %v", err)
	}
	defer pubsubClient.Close()

	// Create notification service
	ns := &NotificationService{
		config:       testConfig,
		db:           testDB,
		pubsubClient: pubsubClient,
	}

	// Test creating a subscriber for daily notifications
	subscriber := pubsubClient.Subscriber(testConfig.DailySubscription)
	if subscriber == nil {
		t.Skip("Cannot create subscriber for testing")
	}

	// Test message processing with a sample event
	testUserID := uuid.New()
	testEvent := models.TrendAnalysisEvent{
		Event:   "trend_analysis",
		JobType: "daily",
		Status:  "success",
		Date:    "2023-10-01",
		SuccessfulUsers: []models.SuccessfulUser{
			{UserID: testUserID, Email: "test@example.com"},
		},
	}

	// Test the handleSuccessEvent method directly (integration test for database)
	err = ns.handleSuccessEvent(ctx, testEvent, DAILY)
	if err != nil {
		t.Fatalf("Expected no error processing success event, got %v", err)
	}

	// Verify notification was created in database
	var notification models.UserNotification
	err = testDB.Where("user_id = ?", testUserID).First(&notification).Error
	if err != nil {
		t.Fatalf("Failed to find notification in database: %v", err)
	}

	if !notification.HasNewDailyTrend {
		t.Error("Expected HasNewDailyTrend to be true")
	}

	if notification.LastDailyTrendDate == nil {
		t.Error("Expected LastDailyTrendDate to be set")
	}

	log.Println("Integration test completed successfully")
}

func setupIntegrationTestDB(t *testing.T) *gorm.DB {
	// Use environment variables for integration test database
	testDSN := os.Getenv("TEST_DATABASE_URL")
	if testDSN == "" {
		testDSN = "postgres://postgres:password@localhost:5432/test_notifications?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Printf("Failed to connect to integration test database: %v", err)
		return nil
	}

	// Auto-migrate
	err = db.AutoMigrate(&models.UserNotification{})
	if err != nil {
		t.Fatalf("Failed to migrate integration test database: %v", err)
	}

	// Clean up existing test data
	db.Where("1 = 1").Delete(&models.UserNotification{})

	return db
}

// Test subscription error handling with database operations
func TestSubscriptionErrorHandling_Integration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test: set RUN_INTEGRATION_TESTS=true to run")
	}

	testDB := setupIntegrationTestDB(t)
	if testDB == nil {
		t.Skip("Database not available for integration testing")
		return
	}

	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, testProjectID, option.WithoutAuthentication())
	if err != nil {
		t.Skipf("PubSub emulator not available: %v", err)
	}
	defer pubsubClient.Close()

	testConfig := &config.AppConfig{
		Profile:           "test",
		ProjectID:         testProjectID,
		DailySubscription: testSubID,
	}

	ns := &NotificationService{
		config:       testConfig,
		db:           testDB,
		pubsubClient: pubsubClient,
	}

	// Test error handling with invalid data
	invalidEvent := models.TrendAnalysisEvent{
		Event:   "", // Missing required field
		JobType: "daily",
		Status:  "success",
		Date:    "2023-10-01",
	}

	err = ns.handleSuccessEvent(ctx, invalidEvent, DAILY)
	if err == nil {
		t.Error("Expected error for invalid event data")
	}

	// Test with valid failure event
	failureEvent := models.TrendAnalysisEvent{
		Event:         "trend_analysis",
		JobType:       "daily",
		Status:        "failure",
		Date:          "2023-10-01",
		FailedUserIDs: []string{"user1", "user2"},
		ErrorSummary:  "Database connection failed",
		NotifyAdmin:   true,
	}

	err = ns.handleFailureEvent(ctx, failureEvent, DAILY)
	if err != nil {
		t.Fatalf("Expected no error for valid failure event, got %v", err)
	}

	log.Println("Error handling integration test completed successfully")
}

// Test notification service operations with database
func TestNotificationService_DatabaseIntegration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test: set RUN_INTEGRATION_TESTS=true to run")
	}

	testDB := setupIntegrationTestDB(t)
	if testDB == nil {
		t.Skip("Database not available for integration testing")
		return
	}

	ctx := context.Background()

	testConfig := &config.AppConfig{
		Profile:   "test",
		ProjectID: testProjectID,
	}

	ns := &NotificationService{
		config: testConfig,
		db:     testDB,
	}

	// Test multiple users with different notification states
	users := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	// Process daily trends for first two users
	dailyEvent := models.TrendAnalysisEvent{
		Event:   "trend_analysis",
		JobType: "daily",
		Status:  "success",
		Date:    "2023-10-01",
		SuccessfulUsers: []models.SuccessfulUser{
			{UserID: users[0], Email: "user1@test.com"},
			{UserID: users[1], Email: "user2@test.com"},
		},
	}

	err := ns.handleSuccessEvent(ctx, dailyEvent, DAILY)
	if err != nil {
		t.Fatalf("Expected no error processing daily event, got %v", err)
	}

	// Process weekly trends for second and third users
	weeklyEvent := models.TrendAnalysisEvent{
		Event:   "trend_analysis",
		JobType: "weekly",
		Status:  "success",
		Date:    "2023-10-01",
		SuccessfulUsers: []models.SuccessfulUser{
			{UserID: users[1], Email: "user2@test.com"},
			{UserID: users[2], Email: "user3@test.com"},
		},
	}

	err = ns.handleSuccessEvent(ctx, weeklyEvent, WEEKLY)
	if err != nil {
		t.Fatalf("Expected no error processing weekly event, got %v", err)
	}

	// Verify notifications for each user
	notifications := make([]models.UserNotification, 3)
	for i, userID := range users {
		err = testDB.Where("user_id = ?", userID).First(&notifications[i]).Error
		if err != nil {
			t.Fatalf("Failed to find notification for user %d: %v", i, err)
		}
	}

	// User 0: Only daily trend
	if !notifications[0].HasNewDailyTrend || notifications[0].HasNewWeeklyTrend {
		t.Error("User 0 should have only daily trend notification")
	}

	// User 1: Both daily and weekly trends
	if !notifications[1].HasNewDailyTrend || !notifications[1].HasNewWeeklyTrend {
		t.Error("User 1 should have both daily and weekly trend notifications")
	}

	// User 2: Only weekly trend
	if notifications[2].HasNewDailyTrend || !notifications[2].HasNewWeeklyTrend {
		t.Error("User 2 should have only weekly trend notification")
	}

	// Test marking notifications as read
	err = ns.MarkNotificationAsRead(users[1], "daily")
	if err != nil {
		t.Fatalf("Expected no error marking daily as read, got %v", err)
	}

	// Verify only daily was marked as read
	var updatedNotification models.UserNotification
	err = testDB.Where("user_id = ?", users[1]).First(&updatedNotification).Error
	if err != nil {
		t.Fatalf("Failed to fetch updated notification: %v", err)
	}

	if updatedNotification.HasNewDailyTrend {
		t.Error("Daily trend should be marked as read")
	}
	if !updatedNotification.HasNewWeeklyTrend {
		t.Error("Weekly trend should still be unread")
	}

	log.Println("Database integration test completed successfully")
}
