//go:build integration

package models

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
	dsn      string
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var db *gorm.DB
var testRepo Repository

func TestMain(m *testing.M) {

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	pool = p

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

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

	resource, err = pool.RunWithOptions(&opts,
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

	log.Println("Waiting 2s to give Postgres container time to boot...")
	time.Sleep(2 * time.Second)

	hostPort := resource.GetPort("5432/tcp")
	os.Setenv("POSTGRES_PORT", hostPort)

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbName, hostPort)
	log.Println("Attempting to connect to database with DSN:", dsn)

	pool.MaxWait = 60 * time.Second

	var tempDB *gorm.DB
	if err = pool.Retry(func() error {
		var openErr error
		tempDB, openErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if openErr != nil {
			log.Printf("Failed to open GORM connection (retrying): %v", openErr)
			return openErr
		}

		sqlDB, sqlErr := tempDB.DB()
		if sqlErr != nil {
			log.Printf("Failed to get underlying SQL DB from GORM (retrying): %v", sqlErr)
			return sqlErr
		}
		pingErr := sqlDB.Ping()
		if pingErr != nil {
			log.Printf("Failed to ping database (retrying): %v", pingErr)
			return pingErr
		}
		log.Println("Successfully pinged database.")
		return nil
	}); err != nil {
		log.Fatalf("Could not establish connection to PostgreSQL database after multiple retries: %s", err)
	}

	db = tempDB

	err = createTables()
	if err != nil {
		log.Fatalf("Could not create tables in the database: %s", err)
	}
	log.Println("✅ Tables initialized successfully.")

	testRepo = &PostgresRepository{DB: db}

	code := m.Run()

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/session.sql")
	if err != nil {
		return fmt.Errorf("error reading session.sql: %w", err)
	}

	if err := db.Exec(string(tableSQL)).Error; err != nil {
		return fmt.Errorf("executing schema session failed: %w", err)
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get database instance: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestPostgresRepository_CreateSession(t *testing.T) {
	repo := &PostgresRepository{DB: db}

	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	// Case 1: Create a session successfully
	t.Run("Create session successfully", func(t *testing.T) {
		session := &Session{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserId:      userID,
			SessionType: "focus",
			StartTime:   time.Now(),
		}

		created, err := repo.CreateSession(session)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		log.Println("created ID:", created.ID)
		if created.ID == uuid.Nil {
			t.Errorf("expected non-nil session ID")
		}
	})

	// Case 2: Fail to create second active session
	t.Run("Fail to create session when one is already active", func(t *testing.T) {
		session := &Session{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserId:      userID,
			SessionType: "focus",
			StartTime:   time.Now().Add(5 * time.Minute),
		}

		_, err := repo.CreateSession(session)
		if err == nil || err.Error() != "user already has an active session — please end it before starting a new one" {
			t.Errorf("expected active session error, got: %v", err)
		}
	})
}

func TestPostgresRepository_StopSession(t *testing.T) {
	repo := &PostgresRepository{DB: db}

	userID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	startTime := time.Now().Add(-1 * time.Hour)

	// Insert active session manually
	session := &Session{
		ID:          uuid.New(),
		UserId:      userID,
		SessionType: "meeting",
		StartTime:   startTime,
		EndTime:     nil,
	}

	if err := db.Create(session).Error; err != nil {
		t.Fatalf("setup failed: could not insert active session: %v", err)
	}

	t.Run("Stop active session", func(t *testing.T) {
		endTime := time.Now()
		sessionDao := &Session{
			UserId:  userID,
			EndTime: &endTime,
		}

		stopped, err := repo.StopSession(sessionDao)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if stopped.EndTime == nil {
			t.Errorf("expected session to have end_time set")
		}
	})

	t.Run("Fail to stop session for user with no active session", func(t *testing.T) {
		sessionDao := &Session{
			UserId:  uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			EndTime: ptrTime(time.Now()),
		}

		_, err := repo.StopSession(sessionDao)
		if err == nil || err.Error() == "" {
			t.Errorf("expected error for no active session, got nil")
		}
	})

	t.Run("UpdateColumn returns error", func(t *testing.T) {
		// Inject an invalid session to simulate failure
		badRepo := &PostgresRepository{DB: db.Session(&gorm.Session{DryRun: true})}
		sessionDao := &Session{
			UserId:  userID,
			EndTime: ptrTime(time.Now()),
		}

		_, err := badRepo.StopSession(sessionDao)
		if err == nil {
			t.Error("expected error due to dry-run DB, got nil")
		}
	})
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
