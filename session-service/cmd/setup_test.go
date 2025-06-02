package main

import (
	"os"
	"productivity-planner/task-service/config"
	"productivity-planner/task-service/models"
	"productivity-planner/task-service/session"

	"testing"

	"github.com/gin-gonic/gin"
)

var appConfig *config.AppConfig
var router *gin.Engine

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	router = gin.New()
	appConfig = config.Load()

	svc := &session.MockSessionService{
		Repo: &models.TestDBRepo{},
	}

	handler := &Handler{Svc: svc}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
