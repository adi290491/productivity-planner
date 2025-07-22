package main

import (
	"os"

	"testing"

	"github.com/adi290491/productivity-planner/user-service/config"
	"github.com/adi290491/productivity-planner/user-service/user"
	"github.com/adi290491/productivity-planner/user-service/utils"
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
