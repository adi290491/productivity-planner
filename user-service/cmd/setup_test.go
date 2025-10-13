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
	// Set required env vars for config.Load()
	os.Setenv("DB_HOSTNAME", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "NxrWXLL7kc")
	os.Setenv("PORT", "1234")
	defer func() {
		os.Unsetenv("DB_HOSTNAME")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
	}()

	// replace shadowing declaration with assignment to package‐level appConfig
	var err error
	appConfig, err = config.Load()
	if err != nil {
		panic("Failed to load config for tests: " + err.Error())
	}

	handler := &Handler{
		Svc:     &user.MockUserService{},
		JwtUtil: utils.JWTUtil{Secret: []byte(appConfig.JWT_SECRET)},
	}

	RegisterEndpoints(router, handler)

	os.Exit(m.Run())
}
