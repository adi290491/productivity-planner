//go:build integration

package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	testRepo *PostgresRepository
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	var err error

	// Create Docker pool
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Start PostgreSQL container
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostPort := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", hostPort)

	// Exponential backoff retry to connect to database
	pool.MaxWait = 30 * time.Second
	if err = pool.Retry(func() error {
		testRepo, err = repository.NewPostgresRepository(dsn)
		if err != nil {
			return err
		}
		return testRepo.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// Load schema
	schema, err := os.ReadFile("../testdata/schema.sql")
	if err != nil {
		log.Fatalf("Could not read schema file: %s", err)
	}

	if _, err := testRepo.DB().Exec(string(schema)); err != nil {
		log.Fatalf("Could not execute schema: %s", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestPostgresRepository_FindSessionsBetweenDates(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		startTime     time.Time
		endTime       time.Time
		expectedCount int
		expectError   bool
	}{
		{
			name:          "user with sessions",
			userID:        "11111111-1111-1111-1111-111111111111",
			startTime:     time.Now().Add(-3 * 24 * time.Hour),
			endTime:       time.Now(),
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "user with no sessions in range",
			userID:        "11111111-1111-1111-1111-111111111111",
			startTime:     time.Now().Add(-10 * 24 * time.Hour),
			endTime:       time.Now().Add(-8 * 24 * time.Hour),
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "non-existent user",
			userID:        "33333333-3333-3333-3333-333333333333",
			startTime:     time.Now().Add(-3 * 24 * time.Hour),
			endTime:       time.Now(),
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "bob's sessions",
			userID:        "22222222-2222-2222-2222-222222222222",
			startTime:     time.Now().Add(-4 * 24 * time.Hour),
			endTime:       time.Now(),
			expectedCount: 3,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := &model.Summary{
				UserId:    tt.userID,
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}

			sessions, err := testRepo.FindSessionsBetweenDates(ctx, summary)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sessions) != tt.expectedCount {
				t.Errorf("expected %d sessions, got %d", tt.expectedCount, len(sessions))
			}

			// Verify session structure
			for _, session := range sessions {
				if session.ID.String() == "" {
					t.Error("session ID is empty")
				}
				if session.UserID.String() != tt.userID {
					t.Errorf("expected user_id %s, got %s", tt.userID, session.UserID.String())
				}
				if session.SessionType == "" {
					t.Error("session type is empty")
				}
				if session.StartTime.IsZero() {
					t.Error("start time is zero")
				}
				if session.EndTime == nil {
					t.Error("end time is nil")
				}
			}
		})
	}
}

func TestPostgresRepository_ConnectionPool(t *testing.T) {
	stats := testRepo.DB().Stats()

	if stats.MaxOpenConnections != 25 {
		t.Errorf("expected MaxOpenConnections to be 25, got %d", stats.MaxOpenConnections)
	}
}

func TestPostgresRepository_Close(t *testing.T) {
	// Create a temporary repository for this test
	hostPort := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", hostPort)

	repo, err := repository.NewPostgresRepository(dsn)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	err = repo.Close()
	if err != nil {
		t.Errorf("expected no error on close, got %v", err)
	}

	// Verify connection is closed by trying to ping
	err = repo.DB().Ping()
	if err == nil {
		t.Error("expected error pinging closed connection, got nil")
	}
}
