package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	origVars := map[string]string{
		"PROFILE":     os.Getenv("PROFILE"),
		"DB_HOSTNAME": os.Getenv("DB_HOSTNAME"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_USERNAME": os.Getenv("DB_USERNAME"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
		"JWT_SECRET":  os.Getenv("JWT_SECRET"),
		"PORT":        os.Getenv("PORT"),
	}

	// Restore env vars after test
	defer func() {
		for key, val := range origVars {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		check   func(*testing.T, *Config)
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"PROFILE":     "test",
				"DB_HOSTNAME": "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USERNAME": "testuser",
				"DB_PASSWORD": "testpass",
				"JWT_SECRET":  "testsecret",
				"PORT":        "8080",
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Profile != "test" {
					t.Errorf("Profile = %v, want test", cfg.Profile)
				}
				if cfg.Port != "8080" {
					t.Errorf("Port = %v, want 8080", cfg.Port)
				}
				if cfg.Database.Host != "localhost" {
					t.Errorf("DB Host = %v, want localhost", cfg.Database.Host)
				}
				if string(cfg.JWT.Secret) != "testsecret" {
					t.Errorf("JWT Secret mismatch")
				}
			},
		},
		{
			name: "missing DB_HOSTNAME",
			envVars: map[string]string{
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USERNAME": "testuser",
				"DB_PASSWORD": "testpass",
				"JWT_SECRET":  "testsecret",
			},
			wantErr: true,
		},
		{
			name: "missing JWT_SECRET",
			envVars: map[string]string{
				"DB_HOSTNAME": "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USERNAME": "testuser",
				"DB_PASSWORD": "testpass",
			},
			wantErr: true,
		},
		{
			name: "default port",
			envVars: map[string]string{
				"DB_HOSTNAME": "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USERNAME": "testuser",
				"DB_PASSWORD": "testpass",
				"JWT_SECRET":  "testsecret",
				// PORT not set
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Port != "8081" {
					t.Errorf("Port = %v, want 8081 (default)", cfg.Port)
				}
			},
		},
		{
			name: "default profile",
			envVars: map[string]string{
				"DB_HOSTNAME": "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USERNAME": "testuser",
				"DB_PASSWORD": "testpass",
				"JWT_SECRET":  "testsecret",
				// PROFILE not set
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Profile != "local" {
					t.Errorf("Profile = %v, want local (default)", cfg.Profile)
				}
			},
		},
		{
			name: "default SSL mode",
			envVars: map[string]string{
				"DB_HOSTNAME": "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USERNAME": "testuser",
				"DB_PASSWORD": "testpass",
				"JWT_SECRET":  "testsecret",
				// DB_SSLMODE not set
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Database.SSLMode != "disable" {
					t.Errorf("SSLMode = %v, want disable (default)", cfg.Database.SSLMode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars
			for key := range origVars {
				os.Unsetenv(key)
			}

			// Set test env vars
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Error("Load() returned nil config")
					return
				}
				if tt.check != nil {
					tt.check(t, cfg)
				}
			}
		})
	}
}

func TestServerConfig(t *testing.T) {
	// Set minimum required env vars
	os.Setenv("DB_HOSTNAME", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "testsecret")

	defer func() {
		os.Unsetenv("DB_HOSTNAME")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Server.ReadTimeout != 10*time.Second {
		t.Errorf("ReadTimeout = %v, want 10s", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout = %v, want 10s", cfg.Server.WriteTimeout)
	}
}

func TestDatabaseDSN(t *testing.T) {
	os.Setenv("DB_HOSTNAME", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("JWT_SECRET", "testsecret")

	defer func() {
		os.Unsetenv("DB_HOSTNAME")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	expectedDSN := "postgres://testuser:testpass@testhost:5433/testdb?sslmode=require"
	if cfg.Database.DSN != expectedDSN {
		t.Errorf("DSN = %v, want %v", cfg.Database.DSN, expectedDSN)
	}
}
