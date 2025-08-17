package config

import (
	"os"
	"testing"
)

// func TestLoad(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		envVars        map[string]string
// 		expectedConfig *AppConfig
// 	}{
// 		{
// 			name: "Valid environment variables",
// 			envVars: map[string]string{
// 				"DB_DSN":               "postgres://user:password@localhost:5432/dbname",
// 				"SESSION_SERVICE_PORT": "8080",
// 			},
// 			expectedConfig: &AppConfig{
// 				DSN:          "postgres://user:password@localhost:5432/dbname",
// 				Port:         "8080",
// 				ReadTimeout:  10 * time.Second,
// 				WriteTimeout: 10 * time.Second,
// 			},
// 		},
// 		{
// 			name:    "Missing environment variables",
// 			envVars: map[string]string{},
// 			expectedConfig: &AppConfig{
// 				DSN:          "",
// 				Port:         "",
// 				ReadTimeout:  10 * time.Second,
// 				WriteTimeout: 10 * time.Second,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			//set env variables
// 			for key, value := range tt.envVars {
// 				os.Setenv(key, value)
// 			}
// 			defer func() {
// 				for key := range tt.envVars {
// 					os.Unsetenv(key)
// 				}
// 			}()

// 			config := Load()

// 			// assert the config values
// 			if config.DSN != tt.expectedConfig.DSN {
// 				t.Errorf("expected DSN %s, got %s", tt.expectedConfig.DSN, config.DSN)
// 			}
// 			if config.Port != tt.expectedConfig.Port {
// 				t.Errorf("expected Port %s, got %s", tt.expectedConfig.Port, config.Port)
// 			}
// 			if config.ReadTimeout != tt.expectedConfig.ReadTimeout {
// 				t.Errorf("expected ReadTimeout %v, got %v", tt.expectedConfig.ReadTimeout, config.ReadTimeout)
// 			}
// 			if config.WriteTimeout != tt.expectedConfig.WriteTimeout {
// 				t.Errorf("expected WriteTimeout %v, got %v", tt.expectedConfig.WriteTimeout, config.WriteTimeout)
// 			}
// 		})
// 	}
// }

func TestLoad(t *testing.T) {
	t.Parallel()
	// Test local profile
	t.Run("local profile loads config.local.yaml", func(t *testing.T) {
		os.Setenv("PROFILE", "local")
		defer os.Unsetenv("PROFILE")

		// Write a temp config.local.yaml
		localConfig := `
profile: local
port: 1234
db:
  username: testuser
  password: testpass
  host: localhost
  port: 5432
  database: testdb
  sslmode: disable
jwt_secret: testsecret
`

		conf, err := os.ReadFile("./config.local.yaml")

		if err != nil {
			t.Fatalf("failed to read original config: %v", err)
		}

		localYaml := []byte(localConfig)
		err = os.WriteFile("./config.local.yaml", localYaml, 0644)
		if err != nil {
			t.Fatalf("failed to write local yaml: %v", err)
		}
		// defer os.Remove("./config.local.yaml")

		config := Load()
		if config.Port != "1234" {
			t.Errorf("expected Port 1234, got %s", config.Port)
		}
		if config.DSN != "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable" {
			t.Errorf("expected DSN, got %s", config.DSN)
		}

		err = os.WriteFile("./config.local.yaml", conf, 0644)
		if err != nil {
			t.Fatalf("failed to write original config: %v", err)
		}
	})

	// Test prod profile with env substitution
	t.Run("prod profile loads config.prod.yaml with env", func(t *testing.T) {
		os.Setenv("PROFILE", "prod")
		defer os.Unsetenv("PROFILE")

		os.Setenv("PORT", "5678")
		os.Setenv("DB_USERNAME", "produser")
		os.Setenv("DB_PASSWORD", "prodpass")
		os.Setenv("DB_HOSTNAME", "prodhost")
		os.Setenv("DB_PORT", "6543")
		os.Setenv("DB_NAME", "proddb")
		os.Setenv("DB_SSLMODE", "require")
		os.Setenv("JWT_SECRET", "prodsecret")
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("DB_USERNAME")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_HOSTNAME")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("DB_SSLMODE")
			os.Unsetenv("JWT_SECRET")
		}()

		// conf, err := os.ReadFile("./config.prod.yaml")

		// if err != nil {
		// 	t.Fatalf("failed to read original config: %v", err)
		// }

		// prodYaml := []byte(`profile: prod\nport: "${PORT}"\ndb:\n  username: ${DB_USERNAME}\n  password: ${DB_PASSWORD}\n  host: ${DB_HOSTNAME}\n  port: ${DB_PORT}\n  database: ${DB_NAME}\n  sslmode: ${DB_SSLMODE}\njwt_secret: ${JWT_SECRET}`)
		// err := os.WriteFile("./config.prod.yaml", prodYaml, 0644)
		// if err != nil {
		// 	t.Fatalf("failed to write prod yaml: %v", err)
		// }
		// defer os.Remove("./config.prod.yaml")

		config := Load()
		if config.Port != "5678" {
			t.Errorf("expected Port 5678, got %s", config.Port)
		}
		if config.DSN != "postgres://produser:prodpass@prodhost:6543/proddb?sslmode=require" {
			t.Errorf("expected DSN, got %s", config.DSN)
		}

		// err = os.WriteFile("./config.prod.yaml", conf, 0644)
		// if err != nil {
		// 	t.Fatalf("failed to write original config: %v", err)
		// }
	})

}
