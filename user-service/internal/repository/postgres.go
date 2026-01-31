package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) UserRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users (id, email, name, password_hash, created_at)
	VALUES ($1, $2, $3, $4, $5)
	`

	slog.Debug("Executing create query", "email", user.Email)
	_, err := r.db.ExecContext(ctx,
		query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.CreatedAt,
	)

	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				slog.Warn("Duplicate email attempted", "email", user.Email)
				return fmt.Errorf("user with email %s already exists", user.Email)
			}
		}
		slog.Error("Failed to create user", "error", err, "email", user.Email)
		return fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("User created successfully", "userId", user.ID, "email", user.Email)
	return nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT id, email, name, password_hash, created_at
	FROM users
	WHERE email = $1
	`

	slog.Debug("Executing GetByEmail query", "email", email)

	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug("User not found", "email", email)
			return nil, fmt.Errorf("user not found")
		}
		slog.Error("Failed to get user by email", "error", err, "email", email)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	slog.Debug("User retrieved successfully", "userId", user.ID, "email", email)
	return &user, nil
}

// GetByIDs retrieves multiple users by their IDs
func (r *PostgresRepository) GetByIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.UserInfo, error) {
	if len(userIDs) == 0 {
		return []model.UserInfo{}, nil
	}

	query := `
		SELECT id, email, name
		FROM users
		WHERE id = ANY($1)
		ORDER BY email
	`

	slog.Debug("Executing GetByIDs query", "count", len(userIDs))

	// Convert []uuid.UUID to pq.Array format
	rows, err := r.db.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		slog.Error("Failed to query users by IDs", "error", err, "count", len(userIDs))
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var userInfos []model.UserInfo
	for rows.Next() {
		var info model.UserInfo
		if err := rows.Scan(&info.ID, &info.Email, &info.Name); err != nil {
			slog.Error("Failed to scan user info", "error", err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userInfos = append(userInfos, info)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating rows", "error", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	slog.Debug("Users retrieved successfully", "found", len(userInfos), "requested", len(userIDs))
	return userInfos, nil
}
