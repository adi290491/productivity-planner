package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/google/uuid"
)

var (
	testDB   *sql.DB
	testRepo UserRepository
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	var err error

	// Create dockertest pool
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Pull and run PostgreSQL container
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
			"listen_addresses='*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Set container to expire in 2 minutes
	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not set expiry: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseURL)

	// Wait for database to be ready
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		testDB, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	log.Println("Database connection established")

	// Create tables
	if err := createTestTables(testDB); err != nil {
		log.Fatalf("Could not create tables: %s", err)
	}

	log.Println("Tables created successfully")

	// Initialize repository
	testRepo = NewPostgresRepository(testDB)

	// Run tests
	code := m.Run()

	// Cleanup
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTestTables(db *sql.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			name TEXT,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`

	_, err := db.Exec(schema)
	return err
}

func cleanupUsers(t *testing.T) {
	_, err := testDB.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		t.Fatalf("Failed to cleanup users: %v", err)
	}
}

func TestPostgresRepository_Create(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	user := &model.User{
		ID:           uuid.New(),
		Email:        "create@example.com",
		Name:         "Create Test",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	err := testRepo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify user was created
	fetchedUser, err := testRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetByEmail() failed: %v", err)
	}

	if fetchedUser.ID != user.ID {
		t.Errorf("ID mismatch: got %v, want %v", fetchedUser.ID, user.ID)
	}
	if fetchedUser.Email != user.Email {
		t.Errorf("Email mismatch: got %v, want %v", fetchedUser.Email, user.Email)
	}
	if fetchedUser.Name != user.Name {
		t.Errorf("Name mismatch: got %v, want %v", fetchedUser.Name, user.Name)
	}
}

func TestPostgresRepository_Create_DuplicateEmail(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	user1 := &model.User{
		ID:           uuid.New(),
		Email:        "duplicate@example.com",
		Name:         "User One",
		PasswordHash: "hash1",
		CreatedAt:    time.Now(),
	}

	err := testRepo.Create(ctx, user1)
	if err != nil {
		t.Fatalf("First Create() failed: %v", err)
	}

	user2 := &model.User{
		ID:           uuid.New(),
		Email:        "duplicate@example.com", // Same email
		Name:         "User Two",
		PasswordHash: "hash2",
		CreatedAt:    time.Now(),
	}

	err = testRepo.Create(ctx, user2)
	if err == nil {
		t.Error("Create() should fail for duplicate email")
	}
}

func TestPostgresRepository_GetByEmail(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// Create test user
	user := &model.User{
		ID:           uuid.New(),
		Email:        "getbyemail@example.com",
		Name:         "Get Test",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	err := testRepo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "getbyemail@example.com",
			wantErr: false,
		},
		{
			name:    "non-existent user",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchedUser, err := testRepo.GetByEmail(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fetchedUser == nil {
				t.Error("GetByEmail() returned nil user")
			}
		})
	}
}

func TestPostgresRepository_GetByIDs(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// Create test users
	user1 := &model.User{
		ID:           uuid.New(),
		Email:        "batch1@example.com",
		Name:         "Batch User 1",
		PasswordHash: "hash1",
		CreatedAt:    time.Now(),
	}
	user2 := &model.User{
		ID:           uuid.New(),
		Email:        "batch2@example.com",
		Name:         "Batch User 2",
		PasswordHash: "hash2",
		CreatedAt:    time.Now(),
	}

	if err := testRepo.Create(ctx, user1); err != nil {
		t.Fatalf("Create user1 failed: %v", err)
	}
	if err := testRepo.Create(ctx, user2); err != nil {
		t.Fatalf("Create user2 failed: %v", err)
	}

	tests := []struct {
		name      string
		userIDs   []uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get both users",
			userIDs:   []uuid.UUID{user1.ID, user2.ID},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "get one user",
			userIDs:   []uuid.UUID{user1.ID},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "empty list",
			userIDs:   []uuid.UUID{},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "non-existent user",
			userIDs:   []uuid.UUID{uuid.New()},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "mixed existing and non-existent",
			userIDs:   []uuid.UUID{user1.ID, uuid.New()},
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := testRepo.GetByIDs(ctx, tt.userIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(users) != tt.wantCount {
				t.Errorf("GetByIDs() count = %v, want %v", len(users), tt.wantCount)
			}
		})
	}
}

func TestPostgresRepository_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond) // Ensure timeout

	user := &model.User{
		ID:           uuid.New(),
		Email:        "timeout@example.com",
		Name:         "Timeout Test",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	// Should fail due to context timeout
	err := testRepo.Create(ctx, user)
	if err == nil {
		t.Error("Create() should fail with context timeout")
	}
}

func TestPostgresRepository_ConcurrentCreates(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// Create multiple users concurrently
	numUsers := 10
	errChan := make(chan error, numUsers)

	for i := 0; i < numUsers; i++ {
		go func(index int) {
			user := &model.User{
				ID:           uuid.New(),
				Email:        fmt.Sprintf("concurrent%d@example.com", index),
				Name:         fmt.Sprintf("Concurrent User %d", index),
				PasswordHash: "hash",
				CreatedAt:    time.Now(),
			}
			errChan <- testRepo.Create(ctx, user)
		}(i)
	}

	// Collect results
	for i := 0; i < numUsers; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("Concurrent Create() failed: %v", err)
		}
	}

	// Verify all users were created
	var count int
	err := testDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if count != numUsers {
		t.Errorf("Expected %d users, got %d", numUsers, count)
	}
}
