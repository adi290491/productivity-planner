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
	appConfig = config.Load()

	svc := trend.NewTrendService(&models.TestDBRepo{})

	handler := &Handler{svc: svc}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
