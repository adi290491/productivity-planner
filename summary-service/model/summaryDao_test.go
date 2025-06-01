//go:build integration

package models

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
	tableSQL, err := os.ReadFile("./testdata/summary.sql")
	if err != nil {
		fmt.Println("error reading summary.sql: %w", err)
		return err
	}

	if err := db.Exec(string(tableSQL)).Error; err != nil {
		return fmt.Errorf("executing schema summary failed: %w", err)
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

func Test_FindAllSessionsBetweenDates(t *testing.T) {
	tests := []struct {
		name           string
		summary        *Summary
		expectedLength int
		expectError    bool
	}{
		{
			name: "has multiple sessions",
			summary: &Summary{
				UserId:    "11111111-1111-1111-1111-111111111111",
				StartTime: time.Now().Add(-3 * 24 * time.Hour), // 3 days ago
				EndTime:   time.Now(),                          // up to current time
			},
			expectedLength: 3,
			expectError:    false,
		},
		{
			name: "has no sessions in time range",
			summary: &Summary{
				UserId:    "11111111-1111-1111-1111-111111111111",
				StartTime: time.Now().Add(-10 * 24 * time.Hour), // 10 days ago
				EndTime:   time.Now().Add(-8 * 24 * time.Hour),  // 8 days ago
			},
			expectedLength: 0, // Expect error or nil
			expectError:    true,
		},
		{
			name: "has one session in narrow window",
			summary: &Summary{
				UserId:    "11111111-1111-1111-1111-111111111111",
				StartTime: time.Now().Add(-25 * time.Hour),
				EndTime:   time.Now().Add(-23 * time.Hour),
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "has sessions for different user",
			summary: &Summary{
				UserId:    "22222222-2222-2222-2222-222222222222",
				StartTime: time.Now().Add(-4 * 24 * time.Hour),
				EndTime:   time.Now(),
			},
			expectedLength: 3,
			expectError:    false,
		},
		{
			name: "invalid user id",
			summary: &Summary{
				UserId:    "33333333-3333-3333-3333-333333333333",
				StartTime: time.Now().Add(-3 * 24 * time.Hour),
				EndTime:   time.Now(),
			},
			expectedLength: 0,
			expectError:    true,
		},
		{
			name: "exact match on session start and end",
			summary: &Summary{
				UserId:    "11111111-1111-1111-1111-111111111111",
				StartTime: time.Now().Add(-48 * time.Hour),
				EndTime:   time.Now().Add(-46 * time.Hour),
			},
			expectedLength: 1,
			expectError:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sessions, err := testRepo.FindAllSessionsBetweenDates(tc.summary)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(sessions) != tc.expectedLength {
				t.Errorf("expected %d sessions, got %d", tc.expectedLength, len(sessions))
			}
		})
	}
}
