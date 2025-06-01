package main

import (
	"os"
	"productivity-planner/trend-service/config"
	models "productivity-planner/trend-service/model"
	"productivity-planner/trend-service/trend"
	"testing"

	"github.com/gin-gonic/gin"
)

var appConfig *config.AppConfig

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	router := gin.New()
	appConfig = config.Load()

	svc := trend.NewTrendService(&models.TestDBRepo{})

	handler := &Handler{svc: svc}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
