package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ProxyToUserService(c *gin.Context) {
	UserServiceUrl := os.Getenv("USER_SERVICE_URL")

	log.Println("User Service URL:", UserServiceUrl+c.Request.URL.Path)

	forward(c, UserServiceUrl+c.Request.URL.Path)
}

func ProxyToSessionService(c *gin.Context) {
	SessionServiceUrl := os.Getenv("SESSION_SERVICE_URL")

	log.Println("Session Service URL:", SessionServiceUrl+c.Request.URL.Path)

	if userId, ok := c.Get("userId"); ok {
		c.Request.Header.Set("X-USER-ID", fmt.Sprintf("%s", userId))
	}

	forward(c, SessionServiceUrl+c.Request.URL.Path)
}

func ProxyToSummaryService(c *gin.Context) {
	summaryServiceUrl := os.Getenv("SUMMARY_SERVICE_URL")

	url := fmt.Sprintf("%s%s?%s", summaryServiceUrl, c.Request.URL.Path, c.Request.URL.RawQuery)
	log.Println("Session Service URL:", url)

	if userId, ok := c.Get("userId"); ok {
		c.Request.Header.Set("X-USER-ID", fmt.Sprintf("%s", userId))
	}

	forward(c, url)
}

func ProxyToTrendService(c *gin.Context) {
	summaryServiceUrl := os.Getenv("TREND_SERVICE_URL")

	url := fmt.Sprintf("%s%s?%s", summaryServiceUrl, c.Request.URL.Path, c.Request.URL.RawQuery)
	log.Println("Trend Service URL:", url)

	if userId, ok := c.Get("userId"); ok {
		c.Request.Header.Set("X-USER-ID", fmt.Sprintf("%s", userId))
	}

	forward(c, url)
}

func forward(c *gin.Context, targetUrl string) {
	var reqBody []byte
	if c.Request.Body != nil {
		reqBody, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	} else {
		reqBody = nil
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, bytes.NewReader(reqBody))

	if err != nil {
		log.Fatal(err)
	}

	req.Header = c.Request.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
		return
	}

	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
