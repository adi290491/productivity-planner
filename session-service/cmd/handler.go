package main

import (
	"fmt"
	"net/http"
	"productivity-planner/task-service/session"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Svc session.SessionServiceInterface
}

func (h *Handler) StartSession(c *gin.Context) {

	// Receive the request body, header
	var sessionRequest session.SessionRequest

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))

	if userId == "" {
		HandleError(c, fmt.Errorf("missing auth token"), http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&sessionRequest); err != nil {
		HandleError(c, fmt.Errorf("invalid request body, %v", err), http.StatusBadRequest)
		return
	}

	if !sessionRequest.SessionType.IsValid() {
		HandleError(c, fmt.Errorf("invalid session type"), http.StatusBadRequest)
		return
	}

	// save the object
	sessionResponse, err := h.Svc.StartSession(sessionRequest, userId)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sessionResponse)
}

func (h *Handler) StopSession(c *gin.Context) {

	var sessionRequest session.SessionRequest

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))

	if userId == "" {
		HandleError(c, fmt.Errorf("missing auth token"), http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&sessionRequest); err != nil {
		HandleError(c, fmt.Errorf("invalid request body, %v", err), http.StatusBadRequest)
		return
	}

	if !sessionRequest.SessionType.IsValid() {
		HandleError(c, fmt.Errorf("invalid session type"), http.StatusBadRequest)
		return
	}

	// save the object
	sessionResponse, err := h.Svc.StopSession(sessionRequest, userId)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sessionResponse)
}
