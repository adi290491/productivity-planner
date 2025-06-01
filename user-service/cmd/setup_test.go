package main

import (
	"os"
	"productivity-planner/user-service/config"
	"productivity-planner/user-service/models"
	"productivity-planner/user-service/user"
	"productivity-planner/user-service/utils"
	"testing"

	"github.com/gin-gonic/gin"
)

var appConfig *config.AppConfig
var router *gin.Engine

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	router = gin.New()
	appConfig = config.Load()

	svc := &user.UserService{Repo: &models.TestDBRepo{}}

	handler := &Handler{
		Svc:     svc,
		JwtUtil: utils.JWTUtil{Secret: []byte("111")},
	}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
