package main

import "github.com/gin-gonic/gin"

type BaseErrorHandler struct {
	Message    string
	StatusCode int
}

func HandleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, BaseErrorHandler{
		Message:    err.Error(),
		StatusCode: statusCode,
	})
}
