package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/adi290491/productivity-planner/daily-trend-analysis-worker/models"
	"github.com/google/uuid"
)

type Publisher struct {
	ProcessingSummary *models.ProcessingSummary
}

func (p *Publisher) Publish(ctx context.Context, app *Application) error {

	if app.PubSubClient == nil {
		return fmt.Errorf("Pub/Sub client is not initialized in the application")
	}

	event := buildTrendAnalysisEvent(p.ProcessingSummary)

	publisher := app.PubSubClient.Publisher(app.PUB_SUB_TOPIC)

	messageData, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("failed to marshall pub/sub event: %w", err)
	}

	result := publisher.Publish(ctx, &pubsub.Message{
		Data: []byte(messageData),
	})

	msgId, err := result.Get(ctx)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Successfully published summary message for job status '%s'; msg ID: %v", event.Status, msgId)

	return nil
}

func buildTrendAnalysisEvent(summary *models.ProcessingSummary) *models.TrendAnalysisEvent {

	dateStr := time.Now().Format("2006-01-02")

	if len(summary.FailedUsers) > 0 {
		var failedIDs []uuid.UUID
		var errorStrings []string
		for _, failedUser := range summary.FailedUsers {
			failedIDs = append(failedIDs, failedUser.UserID)
			errorStrings = append(errorStrings, fmt.Sprintf("UserID %s: %v", failedUser.UserID, failedUser.Errors))
		}

		return &models.TrendAnalysisEvent{
			Event:         "DAILY_TREND_FAILED",
			JobType:       "daily",
			Status:        "failure",
			Date:          dateStr,
			FailedUserIDs: failedIDs,
			ErrorSummary:  strings.Join(errorStrings, "; "),
			NotifyAdmin:   true,
		}
	} else {
		return &models.TrendAnalysisEvent{
			Event:           "DAILY_TREND_COMPLETED",
			JobType:         "daily",
			Status:          "success",
			Date:            dateStr,
			SuccessfulUsers: summary.SuccessfulUsers,
			NotifyAdmin:     false,
		}
	}

}
