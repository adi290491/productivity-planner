package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *AppConfig
	}{
		{
			name: "Valid environment variables",
			envVars: map[string]string{
				"DB_DSN":               "postgres://user:password@localhost:5432/dbname",
				"SUMMARY_SERVICE_PORT": "8080",
			},
			expectedConfig: &AppConfig{
				DSN:          "postgres://user:password@localhost:5432/dbname",
				Port:         "8080",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			},
		},
		{
			name:    "Missing environment variables",
			envVars: map[string]string{},
			expectedConfig: &AppConfig{
				DSN:          "",
				Port:         "",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//set env variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key) //clean up after test
			}

			config := Load()

			// assert the config values
			if config.DSN != tt.expectedConfig.DSN {
				t.Errorf("expected DSN %s, got %s", tt.expectedConfig.DSN, config.DSN)
			}
			if config.Port != tt.expectedConfig.Port {
				t.Errorf("expected Port %s, got %s", tt.expectedConfig.Port, config.Port)
			}
			if config.ReadTimeout != tt.expectedConfig.ReadTimeout {
				t.Errorf("expected ReadTimeout %v, got %v", tt.expectedConfig.ReadTimeout, config.ReadTimeout)
			}
			if config.WriteTimeout != tt.expectedConfig.WriteTimeout {
				t.Errorf("expected WriteTimeout %v, got %v", tt.expectedConfig.WriteTimeout, config.WriteTimeout)
			}
		})
	}
}
