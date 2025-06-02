package main

import (
	"os"
	"productivity-planner/summary-service/config"
	models "productivity-planner/summary-service/model"
	"productivity-planner/summary-service/summary"
	"testing"

	"github.com/gin-gonic/gin"
)

var appConfig *config.AppConfig
var router *gin.Engine

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	router = gin.New()
	appConfig = config.Load()

	svc := &summary.MockSummaryService{
		Repo: &models.TestDBRepo{},
	}

	handler := &Handler{Svc: svc}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
