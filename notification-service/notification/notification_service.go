package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

const (
	DAILY  string = "daily"
	WEEKLY string = "weekly"
)

type NotificationService struct {
	config       *config.AppConfig
	pubsubClient *pubsub.Client

	mu                sync.RWMutex
	messagesProcessed int
	lastProcessed     time.Time

	db *gorm.DB
}

type ProcessingStats struct {
	MessagesProcessed int `json:"messages_processed"`

	LastProcessed time.Time `json:"last_processed"`
}

func NewNotificationService(config *config.AppConfig) (*NotificationService, error) {
	ctx := context.Background()

	var client *pubsub.Client
	var err error
	if config.Profile == "local" {
		credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if credFile != "" {
			log.Println("Using service account key from GOOGLE_APPLICATION_CREDENTIALS")
			client, err = pubsub.NewClient(ctx, config.ProjectID, option.WithCredentialsFile(credFile))
		}
	} else {
		client, err = pubsub.NewClient(ctx, config.ProjectID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return &NotificationService{
		config:       config,
		pubsubClient: client,

		db: config.DB,
	}, nil
}

func (ns *NotificationService) ProcessDailyTrendNotifications(ctx context.Context) error {
	return ns.processNotifications(ctx, DAILY, ns.config.DailySubscription)
}

func (ns *NotificationService) ProcessWeeklyTrendNotifications(ctx context.Context) error {
	return ns.processNotifications(ctx, WEEKLY, ns.config.WeeklySubscription)
}

func (ns *NotificationService) processNotifications(ctx context.Context, trendType string, subscriptionName string) error {
	log.Printf("Starting to process %s trend notifications...", trendType)

	log.Printf("ProjectID: %s", ns.config.ProjectID)
	log.Printf("Using subscription: %s", subscriptionName)
	subscriber := ns.pubsubClient.Subscriber(subscriptionName)

	log.Printf("Subscriber: %+v\n", subscriber)
	subscriber.ReceiveSettings = pubsub.ReceiveSettings{
		MaxExtension:               10 * time.Minute,
		MaxDurationPerAckExtension: 30 * time.Second,
		MinDurationPerAckExtension: 0,
		MaxOutstandingMessages:     1,   // Reduced to prevent overwhelming
		MaxOutstandingBytes:        1e6, // Reduced from 10e6
		NumGoroutines:              1,
	}

	processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	messageCount := 0
	processedMessageIDs := make(map[string]bool)

	err := subscriber.Receive(processCtx, func(ctx context.Context, msg *pubsub.Message) {

		if processedMessageIDs[msg.ID] {
			log.Printf("Skipping duplicate message: %s", msg.ID)
			msg.Ack()
			return
		}

		messageCount++
		processedMessageIDs[msg.ID] = true

		log.Printf("Processing %s trend message %d: ID=%s, Data length: %d bytes",
			trendType, messageCount, msg.ID, len(msg.Data))

		if err := ns.processMessage(ctx, msg, trendType); err != nil {
			log.Printf("Error processing message %s: %v", msg.ID, err)
			msg.Nack()
		} else {
			log.Printf("Successfully processed message %s", msg.ID)
			msg.Ack()
		}
	})

	if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
		return fmt.Errorf("error during message receive: %w", err)
	}

	ns.mu.Lock()
	ns.lastProcessed = time.Now()
	ns.mu.Unlock()

	log.Printf("Completed processing daily trend notifications. Processed %d messages", messageCount)
	return nil
}

func (ns *NotificationService) processMessage(ctx context.Context, msg *pubsub.Message, trendType string) error {
	// Add detailed logging for debugging
	log.Printf("Raw message data: %s", string(msg.Data))

	var event models.TrendAnalysisEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if event.Status == "" {
		return fmt.Errorf("status field is empty or missing")
	}
	if event.JobType == "" {
		return fmt.Errorf("jobType field is empty or missing")
	}

	log.Printf("Processing event: %s, Status: %s, Date: %s", event.Event, event.Status, event.Date)

	// Validate required fields
	if event.Event == "" || event.Status == "" || event.Date == "" {
		log.Printf("Invalid event data - missing required fields. Event: '%s', Status: '%s', Date: '%s'",
			event.Event, event.Status, event.Date)
		return fmt.Errorf("invalid event data: missing required fields")
	}

	ns.mu.Lock()
	ns.messagesProcessed++
	ns.mu.Unlock()

	switch event.Status {
	case "success":
		return ns.handleSuccessEvent(ctx, event, trendType)
	case "failure":
		return ns.handleFailureEvent(ctx, event, trendType)
	default:
		return fmt.Errorf("unknown event status: %s", event.Status)
	}
}

func (ns *NotificationService) handleSuccessEvent(ctx context.Context, event models.TrendAnalysisEvent, trendType string) error {
	log.Printf("Handling %s success event for %d users", trendType, len(event.SuccessfulUsers))

	if len(event.SuccessfulUsers) == 0 {
		log.Println("No successful users to process")
		return nil
	}

	trendDate, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	successCount := 0
	errorCount := 0

	for _, user := range event.SuccessfulUsers {
		userUUID, err := uuid.Parse(user.UserID.String())
		if err != nil {
			log.Printf("Invalid user ID %s: %v", user.UserID, err)
			errorCount++
			continue
		}

		if err := ns.updateNotificationFlag(userUUID, trendType, trendDate, nil); err != nil {
			log.Printf("Failed to update notification for user %s: %v", userUUID, err)
			errorCount++
			continue
		}

		successCount++
		log.Printf("Updated notification flag for user %s (%s)", user.Email, userUUID)
	}

	log.Printf("Processed %d/%d users successfully (%d errors)",
		successCount, len(event.SuccessfulUsers), errorCount)
	return nil
}

func (ns *NotificationService) handleFailureEvent(ctx context.Context, event models.TrendAnalysisEvent, trendType string) error {
	log.Printf("Handling %s failure event: %d failed users", trendType, len(event.FailedUserIDs))
	log.Printf("Error Summary: %s", event.ErrorSummary)

	// Just log for now - in the future, send Slack alerts here
	if event.NotifyAdmin {
		log.Printf("ADMIN ALERT: %s trend processing failed for %d users",
			trendType, len(event.FailedUserIDs))
		log.Printf("Failed User IDs: %v", event.FailedUserIDs)
		log.Printf("Error Details: %s", event.ErrorSummary)
		// TODO: Send Slack notification here
	}

	return nil
}

func (ns *NotificationService) updateNotificationFlag(userID uuid.UUID, trendType string, trendDate time.Time, trendID *uuid.UUID) error {

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	switch trendType {
	case DAILY:
		updates["has_new_daily_trend"] = true
		updates["last_daily_trend_date"] = trendDate
		if trendID != nil {
			updates["last_daily_trend_id"] = trendID
		}
	case WEEKLY:
		updates["has_new_weekly_trend"] = true
		updates["last_weekly_trend_date"] = trendDate
		if trendID != nil {
			updates["last_weekly_trend_id"] = trendID
		}
	}

	result := ns.db.Model(&models.UserNotification{}).
		Where("UserID = ?", userID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update notification: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		notification := &models.UserNotification{
			UserID: userID,
		}

		if trendType == DAILY {
			notification.HasNewDailyTrend = true
			notification.LastDailyTrendDate = &trendDate
			notification.LastDailyTrendID = trendID
		} else if trendType == WEEKLY {
			notification.HasNewWeeklyTrend = true
			notification.LastWeeklyTrendDate = &trendDate
			notification.LastWeeklyTrendID = trendID
		}

		if err := ns.db.Create(notification).Error; err != nil {
			return fmt.Errorf("failed to create notification: %w", err)
		}
	}

	return nil
}

func (ns *NotificationService) GetUserNotification(userID uuid.UUID) (*models.UserNotificationResponse, error) {
	var notification models.UserNotification
	result := ns.db.Where("UserID = ?", userID).First(&notification)

	if result.Error == gorm.ErrRecordNotFound {
		return &models.UserNotificationResponse{
			HasNewDailyTrend:    false,
			HasNewWeeklyTrend:   false,
			LastDailyTrendDate:  nil,
			LastWeeklyTrendDate: nil,
		}, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &models.UserNotificationResponse{
		HasNewDailyTrend:    notification.HasNewDailyTrend,
		HasNewWeeklyTrend:   notification.HasNewWeeklyTrend,
		LastDailyTrendDate:  notification.LastDailyTrendDate,
		LastWeeklyTrendDate: notification.LastWeeklyTrendDate,
	}, nil

}

func (ns *NotificationService) GetStats() ProcessingStats {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ProcessingStats{
		MessagesProcessed: ns.messagesProcessed,
		LastProcessed:     ns.lastProcessed,
	}
}

// MarkNotificationAsRead marks a specific notification type as read for a user
func (ns *NotificationService) MarkNotificationAsRead(userID uuid.UUID, notificationType string) error {
	log.Printf("Marking notification as read for user %s, type: %s", userID, notificationType)

	var updates map[string]interface{}

	switch notificationType {
	case "daily":
		updates = map[string]interface{}{
			"has_new_daily_trend": false,
		}
	case "weekly":
		updates = map[string]interface{}{
			"has_new_weekly_trend": false,
		}
	default:
		return fmt.Errorf("invalid notification type: %s", notificationType)
	}

	result := ns.db.Model(&models.UserNotification{}).
		Where("user_id = ?", userID).
		Updates(updates)

	if result.Error != nil {
		log.Printf("Error marking notification as read: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("No notification found for user %s", userID)
		return fmt.Errorf("notification not found for user")
	}

	log.Printf("Successfully marked %s notification as read for user %s", notificationType, userID)
	return nil
}

func (ns *NotificationService) Close() error {
	if ns.pubsubClient != nil {
		return ns.pubsubClient.Close()
	}

	return nil
}
