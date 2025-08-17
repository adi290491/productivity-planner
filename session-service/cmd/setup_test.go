package main

import (
	"os"

	"testing"

	"github.com/adi290491/productivity-planner/session-service/config"
	"github.com/adi290491/productivity-planner/session-service/models"
	"github.com/adi290491/productivity-planner/session-service/session"
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
