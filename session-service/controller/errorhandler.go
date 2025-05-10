package controller

import "github.com/gin-gonic/gin"

type APIError struct {
	Message    string
	StatusCode int
}

func HandleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, APIError{
		Message:    err.Error(),
		StatusCode: statusCode,
	})
}
