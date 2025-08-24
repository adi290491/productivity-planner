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
	appConfig, _ = config.Load()

	svc := &summary.MockSummaryService{
		Repo: &models.TestDBRepo{},
	}

	handler := &Handler{Svc: svc}

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

/*
create a cloud run deploy command for the newly created summary-service. 
Use these details
region: us-central, 
image tag: us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/summary-service:manual-v1.
env variables:
PROFILE=prod

secrets
DB_HOSTNAME=DB_HOSTNAME version:latest
DB_NAME=DB_NAME version:latest
DB_PASSWORD=DB_PASSWORD version:lates
DB_PORT=DB_PORT version:latest
DB_SSLMODE=DB_SSLMODE versoin:latest
DB_USERNAME=DB_USERNAME version:latest

allow unauthenticated
Cloud SQL Connection: systemic-productivity-planner:us-central1:planner-prod
connect to vcpc from outbound traffic
- directly send traffic to a VCP
*/