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
		var err error
		reqBody, err = io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("Error creating request to %s: %v", targetUrl, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Clone headers but ensure we don't copy problematic headers
	req.Header = c.Request.Header.Clone()

	// Remove Content-Length to avoid conflicts since we're using a new reader
	req.Header.Del("Content-Length")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error forwarding request to %s: %v", targetUrl, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
		return
	}

	defer resp.Body.Close()

	// Copy headers from response, but be selective to avoid conflicts
	for k, v := range resp.Header {
		// Skip headers that Gin handles automatically
		if k == "Content-Length" || k == "Transfer-Encoding" {
			continue
		}
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}

	// Use io.Copy instead of DataFromReader for better reliability
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}
