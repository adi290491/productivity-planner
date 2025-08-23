package main

import (
	"os"

	"testing"

	"github.com/adi290491/productivity-planner/summary-service/config"
	models "github.com/adi290491/productivity-planner/summary-service/model"
	"github.com/adi290491/productivity-planner/summary-service/summary"
	"github.com/gin-gonic/gin"
)

var appConfig *config.AppConfig
var router *gin.Engine

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	router = gin.New()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "NxrWXLL7kc")
	os.Setenv("PORT", "1234")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
	}()
	appConfig, _ = config.Load()

	svc := &summary.MockSummaryService{
		Repo: &models.TestDBRepo{},
	}

	handler := &Handler{Svc: svc}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
