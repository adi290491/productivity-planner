package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleError_SetsCorrectStatusAndBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	err := errors.New("something went wrong")
	status := http.StatusBadRequest

	HandleError(c, err, status)

	if w.Code != status {
		t.Errorf("expected status %d, got %d", status, w.Code)
	}

	expected := `"Message":"something went wrong"`
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("expected response body to contain %q, got %q", expected, w.Body.String())
	}

	expectedStatus := `"StatusCode":400`
	if !strings.Contains(w.Body.String(), expectedStatus) {
		t.Errorf("expected response body to contain %q, got %q", expectedStatus, w.Body.String())
	}
}


func TestAPIError_StructFields_Unique(t *testing.T) {
	apiErr := APIError{
		Message:    "test error",
		StatusCode: 418,
	}
	if apiErr.Message != "test error" {
		t.Errorf("expected Message to be 'test error', got %q", apiErr.Message)
	}
	if apiErr.StatusCode != 418 {
		t.Errorf("expected StatusCode to be 418, got %d", apiErr.StatusCode)
	}
}
