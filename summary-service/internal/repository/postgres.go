package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	slog.Info("Connecting to PostgreSQL database")

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	log.Printf("Database connection successful")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL database")

	return &PostgresRepository{db: db}, nil
}

func (p *PostgresRepository) FindSessionsBetweenDates(ctx context.Context, summary *model.Summary) ([]model.Session, error) {

	query := `
	SELECT id, user_id, session_type, start_time, end_time
	FROM sessions
	WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3 and end_time IS NOT NULL
	ORDER BY start_time ASC
	`

	slog.Debug("Executing query to find sessions",
		"user_id", summary.UserId,
		"start_time", summary.StartTime,
		"end_time", summary.EndTime,
	)

	rows, err := p.db.QueryContext(ctx, query, summary.UserId, summary.StartTime, summary.EndTime)
	if err != nil {
		slog.Error("Failed to execute query",
			"error", err,
			"user_id", summary.UserId,
		)
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}

	defer rows.Close()

	var sessions []model.Session

	for rows.Next() {
		var session model.Session
		var endTime sql.NullTime

		err := rows.Scan(
			&session.ID,
			&session.UserId,
			&session.SessionType,
			&session.StartTime,
			&endTime,
		)

		if err != nil {
			slog.Error("Failed to scan session row",
				"error", err,
				"user_id", summary.UserId,
			)
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		if endTime.Valid {
			session.EndTime = &endTime.Time
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over rows",
			"error", err,
			"user_id", summary.UserId,
		)
		return nil, fmt.Errorf("error iterating over sessions: %w", err)
	}

	if len(sessions) == 0 {
		slog.Debug("No sessions found for user in date range",
			"user_id", summary.UserId,
			"start_time", summary.StartTime,
			"end_time", summary.EndTime,
		)
		return nil, fmt.Errorf("no sessions found for the given day")
	}

	slog.Debug("Sessions retrieved successfully",
		"user_id", summary.UserId,
		"count", len(sessions),
	)

	return sessions, nil

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// var sessions []model.Session

	// err := p.DB.WithContext(ctx).
	// 	Where("user_id = ? AND start_time >= ? AND end_time <= ?", summaryDao.UserId, summaryDao.StartTime, summaryDao.EndTime).
	// 	Find(&sessions).Error

	// if len(sessions) == 0 {
	// 	return nil, fmt.Errorf("no sessions found for the given day")
	// }

	// if err != nil {
	// 	return nil, err
	// }

	// return sessions, nil
}

func (p *PostgresRepository) Close() error {
	slog.Info("Closing database connection")
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (r *PostgresRepository) DB() *sql.DB {
	return r.db
}
