package main

import (
	"os"
	"productivity-planner/user-service/config"
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
	appConfig.JWT_SECRET = "NxrWXLL7kc"

	handler := &Handler{
		Svc:     &user.MockUserService{},
		JwtUtil: utils.JWTUtil{Secret: []byte(appConfig.JWT_SECRET)},
	}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
