package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
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

	client, err := pubsub.NewClient(ctx, config.ProjectID)
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

	// subscriptionName := fmt.Sprintf("/projects/%s/subscriptions/%s", ns.config.ProjectID, ns.config.DailySubscription)

	subscriber := ns.pubsubClient.Subscriber(ns.config.DailySubscription)

	subscriber.ReceiveSettings = pubsub.ReceiveSettings{
		MaxExtension:               10 * time.Minute,
		MaxDurationPerAckExtension: 30 * time.Second,
		MinDurationPerAckExtension: 0,
		MaxOutstandingMessages:     10,
		MaxOutstandingBytes:        10e6,
		NumGoroutines:              1,
	}

	processCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	messageCount := 0

	err := subscriber.Receive(processCtx, func(ctx context.Context, msg *pubsub.Message) {
		messageCount++

		log.Printf("Processing message %d: ID=%s", messageCount, msg.ID)

		if err := ns.processMessage(ctx, msg); err != nil {
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

func (ns *NotificationService) processMessage(ctx context.Context, msg *pubsub.Message) error {
	var event models.TrendAnalysisEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("Processing event: %s, Statsus: %s, Date: %s", event.Event, event.Status, event.Date)

	ns.mu.Lock()
	ns.messagesProcessed++
	ns.mu.Unlock()

	switch event.Status {
	case "success":
		return ns.handleSuccessEvent(ctx, event)
	case "failure":
		return ns.handleFailureEvent(ctx, event)
	default:
		return fmt.Errorf("Unknown event status: %s", event.Status)
	}
}

func (ns *NotificationService) handleSuccessEvent(ctx context.Context, event models.TrendAnalysisEvent) error {
	log.Printf("Handling success event for %d users", len(event.SuccessfulUsers))
	return nil
}

func (ns *NotificationService) handleFailureEvent(ctx context.Context, event models.TrendAnalysisEvent) error {
	log.Printf("Handling failure event: %d failed users", len(event.FailedUserIDs))
	return nil
}

func (ns *NotificationService) GetStats() ProcessingStats {
	ns.mu.RLock()
	defer ns.mu.Unlock()

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
