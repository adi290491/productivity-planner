//go:build integration

package postgres

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "127.0.0.1"
	user     = "postgres"
	password = "postgres"
	dbName   = "productivity_planner_test"
)

var testDB *gorm.DB
var testRepo *SessionRepository

func TestMain(m *testing.M) {
	// Create Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
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
		log.Fatalf("Could not start PostgreSQL container: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	resource.Expire(60)

	log.Println("Waiting for PostgreSQL to be ready...")
	time.Sleep(2 * time.Second)

	hostPort := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbName, hostPort)

	log.Println("Connecting to database:", dsn)

	pool.MaxWait = 60 * time.Second

	// Retry connection
	var db *gorm.DB
	if err = pool.Retry(func() error {
		var openErr error
		db, openErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if openErr != nil {
			log.Printf("Failed to open GORM connection (retrying): %v", openErr)
			return openErr
		}

		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			log.Printf("Failed to get underlying SQL DB (retrying): %v", sqlErr)
			return sqlErr
		}

		pingErr := sqlDB.Ping()
		if pingErr != nil {
			log.Printf("Failed to ping database (retrying): %v", pingErr)
			return pingErr
		}

		log.Println("Successfully connected to database")
		return nil
	}); err != nil {
		log.Fatalf("Could not establish connection to PostgreSQL: %s", err)
	}

	testDB = db

	// Create tables
	err = createTables()
	if err != nil {
		log.Fatalf("Could not create tables: %s", err)
	}
	log.Println("✅ Tables initialized successfully")

	testRepo = NewSessionRepository(testDB)

	code := m.Run()
	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("../../testdata/schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema.sql: %w", err)
	}

	if err := testDB.Exec(string(tableSQL)).Error; err != nil {
		return fmt.Errorf("executing schema failed: %w", err)
	}

	return nil
}

func TestPingDB(t *testing.T) {
	sqlDB, err := testDB.DB()
	if err != nil {
		t.Fatalf("Failed to get database instance: %v", err)
	}

	err = sqlDB.Ping()
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
	testDB.Exec("DELETE FROM sessions WHERE user_id = ?", userID.String())
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
	testDB.Exec("DELETE FROM sessions WHERE user_id = ?", userID.String())
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
	testDB.Exec("DELETE FROM sessions WHERE user_id = ?", userID.String())
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
