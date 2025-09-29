package models

import "github.com/google/uuid"

type TrendAnalysisEvent struct {
	Event           string           `json:"event"`
	JobType         string           `json:"jobType"`
	Status          string           `json:"status"`
	Date            string           `json:"date"`
	SuccessfulUsers []SuccessfulUser `json:"successfulUsers,omitempty"`
	FailedUserIDs   []string         `json:"failedUserIds,omitempty"`
	ErrorSummary    string           `json:"errorSummary,omitempty"`
	NotifyAdmin     bool             `json:"notifyAdmin"`
}

type SuccessfulUser struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

type UserFailureInfo struct {
	UserID uuid.UUID
	Errors error
}
