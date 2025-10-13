package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterEndpoints_Routes(t *testing.T) {
	expectedRoutes := []struct {
		route  string
		method string
	}{
		{"/users/signup", "POST"},
		{"/users/login", "POST"},
		{"/users/batch", "POST"},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	var handler *Handler
	RegisterEndpoints(router, handler)

	for _, expectedRoute := range expectedRoutes {
		if !routeExists(expectedRoute.route, expectedRoute.method, router) {
			t.Errorf("Route %s with method %s not registered", expectedRoute.route, expectedRoute.method)
		}
	}
}

func TestRegisterEndpoints_RouteCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	var handler *Handler
	RegisterEndpoints(router, handler)

	routes := router.Routes()
	expectedRouteCount := 3

	if len(routes) != expectedRouteCount {
		t.Errorf("Expected %d routes, got %d", expectedRouteCount, len(routes))
	}
}

func TestRegisterEndpoints_NilHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Should not panic with nil handler
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterEndpoints panicked with nil handler: %v", r)
		}
	}()

	RegisterEndpoints(router, nil)
}

func routeExists(testRoute, testMethod string, routes *gin.Engine) bool {
	for _, route := range routes.Routes() {
		if route.Path == testRoute && route.Method == testMethod {
			return true
		}
	}
	return false
}
