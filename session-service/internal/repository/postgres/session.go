package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SessionRepository implements repository.SessionRepository using PostgreSQL
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new PostgreSQL session repository
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// NewDB creates a new database connection
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	return db, nil
}

// dbSession represents the database schema for sessions table
type dbSession struct {
	ID          string     `gorm:"column:id;primaryKey"`
	UserID      string     `gorm:"column:user_id"`
	SessionType string     `gorm:"column:session_type"`
	StartTime   time.Time  `gorm:"column:start_time"`
	EndTime     *time.Time `gorm:"column:end_time"`
}

func (dbSession) TableName() string {
	return "sessions"
}

// Create creates a new session in the database
func (r *SessionRepository) Create(session *model.Session) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var existing dbSession

	// Check if active session exists
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time IS NULL", session.UserID.String()).
		First(&existing).Error

	// If an active session exists, return error
	if err == nil {
		return nil, errors.New("user already has an active session — please end it before starting a new one")
	}

	// If error is not "record not found", return the error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking for active session: %w", err)
	}

	// Create new session
	dbSess := &dbSession{
		ID:          session.ID.String(),
		UserID:      session.UserID.String(),
		SessionType: session.SessionType.String(),
		StartTime:   session.StartTime,
		EndTime:     session.EndTime,
	}

	err = r.db.WithContext(ctx).Create(dbSess).Error
	if err != nil {
		return nil, fmt.Errorf("session creation failed: %w", err)
	}

	return session, nil
}

// Stop updates an active session with an end time
func (r *SessionRepository) Stop(session *model.Session) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var existing dbSession

	// Find active session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time IS NULL", session.UserID.String()).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no active session found: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	// Update end time
	existing.EndTime = session.EndTime

	result := r.db.WithContext(ctx).
		Model(&existing).
		UpdateColumn("end_time", existing.EndTime)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update session: %w", result.Error)
	}

	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("unexpected number of rows affected: %d", result.RowsAffected)
	}

	// Convert back to model
	endTime := *existing.EndTime
	stoppedSession := &model.Session{
		ID:          session.UserID, // Use the ID from the found session
		UserID:      session.UserID,
		SessionType: model.SessionType(existing.SessionType),
		StartTime:   existing.StartTime,
		EndTime:     &endTime,
	}

	log.Printf("Session stopped successfully for user %s", session.UserID)
	return stoppedSession, nil
}
