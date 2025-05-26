package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterEndpoints_Routes(t *testing.T) {
	// Your test implementation here
	var registered = []struct {
		route  string
		method string
	}{
		{"/trend/daily", "GET"},
		{"/trend/weekly", "GET"},
	}

	router := gin.New()
	var handler *Handler
	RegisterEndpoints(router, handler)
	for _, r := range registered {
		if !routeExists(r.route, r.method, router) {
			t.Errorf("Route %s with method %s not registered", r.route, r.method)
		}
	}
}

func routeExists(testRoute, testMethod string, routes *gin.Engine) bool {
	for _, route := range routes.Routes() {
		if route.Path == testRoute && route.Method == testMethod {
			return true
		}
	}
	return false
}
