package main

import (
	"os"

	"testing"

	"github.com/adi290491/productivity-planner/trend-service/config"
	models "github.com/adi290491/productivity-planner/trend-service/model"
	"github.com/adi290491/productivity-planner/trend-service/trend"
	"github.com/gin-gonic/gin"
)

var appConfig *config.AppConfig

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	router := gin.New()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "NxrWXLL7kc")
	os.Setenv("PORT", "1234")
	appConfig, _ = config.Load()

	svc := trend.NewTrendService(&models.TestDBRepo{})

	handler := &Handler{svc: svc}

	RegisterEndpoints(router, handler)

	cleanup := func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
	}
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
	os.Exit(m.Run())
}
