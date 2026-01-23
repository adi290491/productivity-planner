package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
)

// SessionRepository implements repository.SessionRepository using PostgreSQL
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new PostgreSQL session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// NewDB creates a new database connection
func NewDB(dsn string) (*sql.DB, error) {
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	slog.Info("Connecting to database", "dsn", dsn)
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
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

	// var existing dbSession

	// // Check if active session exists
	// err := r.db.WithContext(ctx).
	// 	Where("user_id = ? AND end_time IS NULL", session.UserID.String()).
	// 	First(&existing).Error

	// // If an active session exists, return error
	// if err == nil {
	// 	return nil, errors.New("user already has an active session — please end it before starting a new one")
	// }

	// // If error is not "record not found", return the error
	// if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, fmt.Errorf("error checking for active session: %w", err)
	// }

	// // Create new session
	// dbSess := &dbSession{
	// 	ID:          session.ID.String(),
	// 	UserID:      session.UserID.String(),
	// 	SessionType: session.SessionType.String(),
	// 	StartTime:   session.StartTime,
	// 	EndTime:     session.EndTime,
	// }

	// err = r.db.WithContext(ctx).Create(dbSess).Error
	// if err != nil {
	// 	return nil, fmt.Errorf("session creation failed: %w", err)
	// }

	// return session, nil

	// Check if active session exists
	var existingID string
	checkQuery := `
	SELECT id
	FROM sessions
	WHERE user_id = $1 AND end_time IS NULL
	LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, checkQuery, session.UserID).Scan(&existingID)

	if err == nil {
		// Active session found
		return nil, errors.New("user already has an active session — please end it before starting a new one")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error checking for active session: %w", err)
	}

	// insert new session
	insertQuery := `
	INSERT INTO sessions (id, user_id, session_type, start_time, end_time)
	Values ($1, $2, $3, $4, $5)
	Returning id
	`

	var returnedId string
	err = r.db.QueryRowContext(ctx,
		insertQuery,
		session.ID,
		session.UserID,
		session.SessionType.String(),
		session.StartTime,
		session.EndTime,
	).Scan(&returnedId)

	if err != nil {
		slog.Error("Failed to insert session", "error", err) // ← Added logging
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	slog.Info("Session created successfully",
		"session_id", session.ID,
		"user_id", session.UserID,
		"type", session.SessionType,
	)
	return session, nil
}

// Stop updates an active session with an end time
func (r *SessionRepository) Stop(session *model.Session) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// var existing dbSession

	// // Find active session
	// err := r.db.WithContext(ctx).
	// 	Where("user_id = ? AND end_time IS NULL", session.UserID.String()).
	// 	First(&existing).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, fmt.Errorf("no active session found: %w", err)
	// }

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query session: %w", err)
	// }

	// // Update end time
	// existing.EndTime = session.EndTime

	// result := r.db.WithContext(ctx).
	// 	Model(&existing).
	// 	UpdateColumn("end_time", existing.EndTime)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("failed to update session: %w", result.Error)
	// }

	// if result.RowsAffected != 1 {
	// 	return nil, fmt.Errorf("unexpected number of rows affected: %d", result.RowsAffected)
	// }

	// // Convert back to model
	// endTime := *existing.EndTime
	// stoppedSession := &model.Session{
	// 	ID:          session.UserID, // Use the ID from the found session
	// 	UserID:      session.UserID,
	// 	SessionType: model.SessionType(existing.SessionType),
	// 	StartTime:   existing.StartTime,
	// 	EndTime:     &endTime,
	// }

	// log.Printf("Session stopped successfully for user %s", session.UserID)
	// return stoppedSession, nil
	updateQuery := `
	UPDATE sessions
	SET end_time = $1
	WHERE user_id = $2 AND end_time IS NULL
	RETURNING id, session_type, start_time, end_time
	`

	var (
		id          string
		sessionType string
		startTime   time.Time
		endTime     time.Time
	)

	err := r.db.QueryRowContext(
		ctx,
		updateQuery,
		session.EndTime,
		session.UserID,
	).Scan(&id, &sessionType, &startTime, &endTime)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no active session found for user")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stop session: %w", err)
	}

	stoppedSession := &model.Session{
		ID:          session.UserID, // Use the ID from the found session
		UserID:      session.UserID,
		SessionType: model.SessionType(sessionType),
		StartTime:   startTime,
		EndTime:     &endTime,
	}

	slog.Info("Session stopped successfully",
		"session_id", id,
		"user_id", session.UserID,
		"duration_seconds", endTime.Sub(startTime).Seconds(),
	)
	return stoppedSession, nil
}

func (r *SessionRepository) Close() error {
	return r.db.Close()
}
