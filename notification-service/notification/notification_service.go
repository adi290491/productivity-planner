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
	"google.golang.org/api/option"
)

type NotificationService struct {
	config       *config.AppConfig
	pubsubClient *pubsub.Client

	mu                sync.RWMutex
	messagesProcessed int
	emailsSent        int
	lastProcessed     time.Time
}

type ProcessingStats struct {
	MessagesProcessed int       `json:"messages_processed"`
	EmailsSent        int       `json:"emails_sent"`
	LastProcessed     time.Time `json:"last_processed"`
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
	}, nil
}

func (ns *NotificationService) ProcessDailyTrendNotifications(ctx context.Context) error {
	log.Println("Starting to process daily trend notifications...")

	// subscriptionName := fmt.Sprintf("projects/%s/subscriptions/%s", ns.config.ProjectID, ns.config.DailySubscription)
	// log.Println("Subscription:", subscriptionName)
	log.Printf("ProjectID: %s", ns.config.ProjectID)
	log.Printf("DailySubscription: %s", ns.config.DailySubscription)
	subscriber := ns.pubsubClient.Subscriber(ns.config.DailySubscription)

	log.Printf("Subscriber: %+v\n", subscriber)
	subscriber.ReceiveSettings = pubsub.ReceiveSettings{
		MaxExtension:               10 * time.Minute,
		MaxDurationPerAckExtension: 30 * time.Second,
		MinDurationPerAckExtension: 0,
		MaxOutstandingMessages:     5,   // Reduced to prevent overwhelming
		MaxOutstandingBytes:        1e6, // Reduced from 10e6
		NumGoroutines:              1,
	}

	processCtx, cancel := context.WithTimeout(ctx, 2*time.Minute) // Reduced from 5 minutes
	defer cancel()

	messageCount := 0
	maxMessages := 10 // Limit the number of messages to process

	err := subscriber.Receive(processCtx, func(ctx context.Context, msg *pubsub.Message) {
		messageCount++

		log.Printf("Processing message %d: ID=%s (max: %d)", messageCount, msg.ID, maxMessages)
		log.Printf("Message delivery attempt: %d", msg.DeliveryAttempt)

		// Stop processing if we've reached the max message limit
		if messageCount > maxMessages {
			log.Printf("Reached maximum message limit (%d), stopping processing", maxMessages)
			cancel() // Cancel the context to stop receiving more messages
			return
		}

		if err := ns.processMessage(ctx, msg); err != nil {
			log.Printf("Error processing message %s: %v", msg.ID, err)

			// Avoid infinite redelivery - ack messages that fail after multiple attempts
			if msg.DeliveryAttempt != nil && *msg.DeliveryAttempt > 3 {
				log.Printf("Message %s failed after %d attempts, acknowledging to prevent infinite redelivery",
					msg.ID, *msg.DeliveryAttempt)
				msg.Ack() // Acknowledge to prevent redelivery
			} else {
				msg.Nack() // Allow retry for the first few attempts
			}
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

func (ns *NotificationService) processMessage(ctx context.Context, msg *pubsub.Message) error {
	// Add detailed logging for debugging
	log.Printf("Raw message data: %s", string(msg.Data))
	log.Printf("Message attributes: %+v", msg.Attributes)

	// Check if this is a test message and skip it
	var testCheck map[string]interface{}
	if err := json.Unmarshal(msg.Data, &testCheck); err == nil {
		if testValue, exists := testCheck["test"]; exists {
			log.Printf("Skipping test message: %v", testValue)
			return nil // Return nil to acknowledge and skip test messages
		}
	}

	var event models.TrendAnalysisEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("JSON unmarshaling failed. Raw data: %s", string(msg.Data))
		return fmt.Errorf("failed to unmarshal message: %w", err)
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
		return ns.handleSuccessEvent(ctx, event)
	case "failure":
		return ns.handleFailureEvent(ctx, event)
	default:
		return fmt.Errorf("unknown event status: %s", event.Status)
	}
}

func (ns *NotificationService) handleSuccessEvent(ctx context.Context, event models.TrendAnalysisEvent) error {
	log.Printf("Handling success event for %d users", len(event.SuccessfulUsers))

	if len(event.SuccessfulUsers) == 0 {
		log.Println("No successful users to notify")
		return nil
	}

	// Send success emails to users in batches
	batchSize := 10 // Process 10 emails at a time
	emailsSent := 0

	for i := 0; i < len(event.SuccessfulUsers); i += batchSize {
		end := i + batchSize
		if end > len(event.SuccessfulUsers) {
			end = len(event.SuccessfulUsers)
		}

		batch := event.SuccessfulUsers[i:end]
		log.Printf("Batch: %+v", batch)
		// for _, user := range batch {
		// 	if err := ns.emailService.SendSuccessEmail(user, event); err != nil {
		// 		log.Printf("Failed to send success email to %s (%s): %v", user.Email, user.UUID, err)
		// 		// Continue with other emails even if one fails
		// 	} else {
		// 		emailsSent++
		// 		log.Printf("Sent success email to %s (%s)", user.Email, user.UUID)
		// 	}
		// }

		// Small delay between batches to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	ns.mu.Lock()
	ns.emailsSent += emailsSent
	ns.mu.Unlock()

	log.Printf("Successfully sent %d success emails out of %d users", emailsSent, len(event.SuccessfulUsers))
	return nil
}

func (ns *NotificationService) handleFailureEvent(ctx context.Context, event models.TrendAnalysisEvent) error {
	log.Printf("Handling failure event: %d failed users", len(event.FailedUserIDs))

	if !event.NotifyAdmin {
		log.Println("Admin notification not required for this failure")
		return nil
	}

	// // Send failure alert to admin
	// if err := ns.emailService.SendFailureAlert(event); err != nil {
	// 	return fmt.Errorf("failed to send failure alert to admin: %w", err)
	// }

	ns.mu.Lock()
	ns.emailsSent++
	ns.mu.Unlock()

	log.Printf("Sent failure alert to admin for %d failed users", len(event.FailedUserIDs))
	return nil
}

func (ns *NotificationService) GetStats() ProcessingStats {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ProcessingStats{
		MessagesProcessed: ns.messagesProcessed,
		LastProcessed:     ns.lastProcessed,
		EmailsSent:        ns.emailsSent,
	}
}

func (ns *NotificationService) Close() error {
	if ns.pubsubClient != nil {
		return ns.pubsubClient.Close()
	}

	return nil
}
