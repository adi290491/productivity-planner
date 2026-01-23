//go:build integration

package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "127.0.0.1"
	user     = "postgres"
	password = "postgres"
	dbName   = "productivity_planner_test"
)

var testDB *sql.DB
var testRepo *SessionRepository

func TestMain(m *testing.M) {
	// Create Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		slog.Error("Could not connect to Docker", "error", err)
		os.Exit(1)
	}

	err = pool.Client.Ping()
	if err != nil {
		slog.Error("Could not connect to Docker", "error", err)
		os.Exit(1)
	}

	// Start PostgreSQL container
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
			"listen_addresses= '*'",
		},
	}

	resource, err := pool.RunWithOptions(&opts,
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
	if err != nil {
		slog.Error("Could not start PostgreSQL container", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			slog.Error("Could not purge resource", "error", err)
		}
	}()

	resource.Expire(60)

	slog.Info("Waiting for PostgreSQL to be ready...")
	time.Sleep(2 * time.Second)

	hostPort := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, hostPort, dbName)

	slog.Info("Connecting to database", "dsn", dsn)

	pool.MaxWait = 60 * time.Second

	// Retry connection
	var db *sql.DB
	if err = pool.Retry(func() error {
		var openErr error
		db, openErr = sql.Open("pgx", dsn)
		if openErr != nil {
			slog.Warn("Failed to open connection (retrying)", "error", openErr)
			return openErr
		}

		pingErr := db.Ping()
		if pingErr != nil {
			slog.Warn("Failed to ping database (retrying)", "error", pingErr)
			return pingErr
		}

		slog.Info("Successfully connected to database")
		return nil
	}); err != nil {
		slog.Error("Could not establish connection to PostgreSQL", "error", err)
		os.Exit(1)
	}

	testDB = db

	// Create tables
	err = createTables()
	if err != nil {
		slog.Error("Could not create tables", "error", err)
		os.Exit(1)
	}
	slog.Info("Tables initialized successfully")

	testRepo = NewSessionRepository(testDB)

	code := m.Run()
	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("../../../testdata/schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema.sql: %w", err)
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		return fmt.Errorf("executing schema failed: %w", err)
	}

	return nil
}

func TestPingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestSessionRepository_Create_Success(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	session := &model.Session{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		StartTime:   time.Now().UTC(),
		EndTime:     nil,
	}

	created, err := testRepo.Create(session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Errorf("expected non-nil session ID")
	}

	// Clean up
	_, _ = testDB.Exec("DELETE FROM sessions WHERE user_id = $1", userID)
}

func TestSessionRepository_Create_AlreadyActive(t *testing.T) {
	userID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	// Create first session
	session1 := &model.Session{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		StartTime:   time.Now().UTC(),
		EndTime:     nil,
	}

	_, err := testRepo.Create(session1)
	if err != nil {
		t.Fatalf("failed to create first session: %v", err)
	}

	// Try to create second active session
	session2 := &model.Session{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: model.SessionTypeMeeting,
		StartTime:   time.Now().UTC(),
		EndTime:     nil,
	}

	_, err = testRepo.Create(session2)
	if err == nil {
		t.Errorf("expected error for already active session, got nil")
	}

	// Clean up
	_, _ = testDB.Exec("DELETE FROM sessions WHERE user_id = $1", userID)
}

func TestSessionRepository_Stop_Success(t *testing.T) {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	// Create active session
	session := &model.Session{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		StartTime:   time.Now().UTC().Add(-30 * time.Minute),
		EndTime:     nil,
	}

	_, err := testRepo.Create(session)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Stop session
	endTime := time.Now().UTC()
	stopSession := &model.Session{
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		EndTime:     &endTime,
	}

	stopped, err := testRepo.Stop(stopSession)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stopped.EndTime == nil {
		t.Errorf("expected session to have end_time set")
	}

	// Clean up
	_, _ = testDB.Exec("DELETE FROM sessions WHERE user_id = $1", userID)
}

func TestSessionRepository_Stop_NoActiveSession(t *testing.T) {
	userID := uuid.New()

	endTime := time.Now().UTC()
	stopSession := &model.Session{
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		EndTime:     &endTime,
	}

	_, err := testRepo.Stop(stopSession)
	if err == nil {
		t.Errorf("expected error for no active session, got nil")
	}
}
