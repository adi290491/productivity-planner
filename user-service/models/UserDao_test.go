//go:build integration

package models

import (
	"errors"
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
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println("error reading init.sql: %w", err)
		return err
	}

	if err := db.Exec(string(tableSQL)).Error; err != nil {
		return fmt.Errorf("executing schema init failed: %w", err)
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

func Test_CreateUser(t *testing.T) {
	testUser := &User{
		ID:           uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Name:         "Test User",
		Email:        "create@example.com",
		PasswordHash: "1a2b3c",
	}

	result, err := testRepo.CreateUser(testUser)

	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// For UUIDs, check if ID is not uuid.Nil
	if result.ID == uuid.Nil {
		t.Errorf("Expected created user to have a non-nil UUID, got %v", result.ID)
	}
	if result.Email != testUser.Email {
		t.Errorf("Expected email %s, got %s", testUser.Email, result.Email)
	}
	if result.PasswordHash != testUser.PasswordHash {
		t.Errorf("Expected password hash %s, got %s", testUser.PasswordHash, result.PasswordHash)
	}

	// Verify the user exists in the database by fetching it
	fetchedUser, err := testRepo.GetUser(&User{Email: testUser.Email})
	if err != nil {
		t.Fatalf("GetUser failed after CreateUser: %v", err)
	}
	if fetchedUser.Email != testUser.Email {
		t.Errorf("Fetched user email mismatch: expected %s, got %s", testUser.Email, fetchedUser.Email)
	}
	if fetchedUser.Name != testUser.Name {
		t.Errorf("Fetched user name mismatch: expected %s, got %s", testUser.Name, fetchedUser.Name)
	}
	if fetchedUser.ID == uuid.Nil {
		t.Errorf("Expected fetched user to have a non-nil UUID, got %v", fetchedUser.ID)
	}
	if fetchedUser.PasswordHash != testUser.PasswordHash {
		t.Errorf("Fetched user password hash mismatch: expected %s, got %s", testUser.PasswordHash, fetchedUser.PasswordHash)
	}
}

// TestGetUser tests the GetUser method of PostgresRepository.
func TestGetUser(t *testing.T) {
	// Ensure a clean state before each test
	t.Cleanup(func() {
		if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
			t.Fatalf("Failed to truncate users table: %v", err)
		}
	})
	existingUser := &User{
		Email: "alice@example.com",
		Name:  "Alice",
		ID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	}

	fetchedUser, err := testRepo.GetUser(existingUser)
	if err != nil {
		t.Fatalf("GetUser failed for existing user: %v", err)
	}

	if fetchedUser.Email != existingUser.Email {
		t.Errorf("Expected fetched user email %s, got %s", existingUser.Email, fetchedUser.Email)
	}
	if fetchedUser.Name != existingUser.Name {
		t.Errorf("Expected fetched user name %s, got %s", existingUser.Name, fetchedUser.Name)
	}
	if fetchedUser.ID == uuid.Nil {
		t.Errorf("Expected fetched user to have a non-nil UUID, got %v", fetchedUser.ID)
	}

	// Test case 2: User does not exist
	nonExistentUserEmail := "nonexistent@example.com"
	_, err = testRepo.GetUser(&User{Email: nonExistentUserEmail})
	if err == nil {
		t.Error("Expected GetUser to return an error for non-existent user, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected error to be gorm.ErrRecordNotFound, got %v", err)
	}
}
