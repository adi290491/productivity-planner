//go:build integration

package main

import (
	"fmt"
	"log"
	"os"
	"productivity-planner/trend-analysis-worker/weekly-aggregates/models"

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
var testRepo *PostgresRepository

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
	tableSQL, err := os.ReadFile("../testdata/worker.sql")
	if err != nil {
		return fmt.Errorf("error reading worker.sql: %w", err)
	}

	if err := db.Exec(string(tableSQL)).Error; err != nil {
		return fmt.Errorf("executing schema worker failed: %w", err)
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

func TestFetchWeeklyTrend(t *testing.T) {
	// repo := &PostgresRepository{DB: db}

	testRepo.FetchWeeklyTrend()

	var trends []models.UserWeeklyTrend
	if err := db.Find(&trends).Error; err != nil {
		t.Fatalf("Failed to fetch weekly trends: %v", err)
	}

	if len(trends) < 2 {
		t.Fatalf("Expected at least 2 weekly trend rows (Alice and Bob), got %d", len(trends))
	}

	foundAlice := false
	foundBob := false

	for _, trend := range trends {
		if trend.UserId.String() == "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa" {
			foundAlice = true
			if trend.FocusMinutes <= 0 || trend.MeetingMinutes <= 0 || trend.BreakMinutes <= 0 {
				t.Errorf("Expected non-zero values for Alice. Got: %+v", trend)
			}
		}
		if trend.UserId.String() == "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb" {
			foundBob = true
			if trend.FocusMinutes <= 0 || trend.MeetingMinutes <= 0 || trend.BreakMinutes <= 0 {
				t.Errorf("Expected non-zero values for Bob. Got: %+v", trend)
			}
		}
	}

	if !foundAlice {
		t.Error("Did not find weekly trend row for Alice")
	}
	if !foundBob {
		t.Error("Did not find weekly trend row for Bob")
	}
}
