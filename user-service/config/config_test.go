package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
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
		if config.JWT_SECRET != "testsecret" {
			t.Errorf("expected JWT_SECRET testsecret, got %s", config.JWT_SECRET)
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
		if config.JWT_SECRET != "prodsecret" {
			t.Errorf("expected JWT_SECRET prodsecret, got %s", config.JWT_SECRET)
		}

		// err = os.WriteFile("./config.prod.yaml", conf, 0644)
		// if err != nil {
		// 	t.Fatalf("failed to write original config: %v", err)
		// }
	})

}

// func TestLoad_ProdProfile(t *testing.T) {
// 	prodConfig := `
// profile: prod
// port: "${PORT}"
// db:
//   username: "${DB_USER}"
//   password: "${DB_PASSWORD}"
//   host: "${DB_HOST}"
//   port: "${DB_PORT}"
//   database: "${DB_NAME}"
//   sslmode: "${DB_SSLMODE}"
// jwt_secret: "${JWT_SECRET}"
// `
// 	path := writeTempConfig(prodConfig)
// 	defer os.Remove(path)

// 	// Simulate secret-injected env vars
// 	t.Setenv("PROFILE", "prod")
// 	t.Setenv("CONFIG_PATH_PROD", path)
// 	t.Setenv("PORT", "9090")
// 	t.Setenv("DB_USER", "prod_user")
// 	t.Setenv("DB_PASSWORD", "prod_pass")
// 	t.Setenv("DB_HOST", "cloudsql")
// 	t.Setenv("DB_PORT", "5432")
// 	t.Setenv("DB_NAME", "prod_db")
// 	t.Setenv("DB_SSLMODE", "disable")
// 	t.Setenv("JWT_SECRET", "prod_secret")

// 	cfg := Load()

// 	assert.Equal(t, "postgres://prod_user:prod_pass@cloudsql:5432/prod_db?sslmode=disable", cfg.DSN)
// 	assert.Equal(t, "prod_secret", cfg.JWT_SECRET)
// 	assert.Equal(t, "9090", cfg.Port)
// 	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
// 	assert.Equal(t, 10*time.Second, cfg.WriteTimeout)
// 	assert.Nil(t, cfg.DB)
// }

// func writeTempConfig(content string) string {
// 	tmpFile, err := os.CreateTemp("", "config-*.yaml")
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = tmpFile.WriteString(content)
// 	if err != nil {
// 		panic(err)
// 	}
// 	tmpFile.Close()
// 	return tmpFile.Name()
// }
