package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
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

	// --- START: MODIFICATION ---
	// Create an authenticated client that adds a Google-signed identity token
	// to outbound requests.
	audience := extractBaseURL(targetUrl)
	ctx := context.Background()

	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		log.Printf("Error creating authenticated client for audience %s: %v", audience, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create internal authenticated client"})
		return
	}
	// --- END: MODIFICATION ---

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

// extractBaseURL extracts the base URL (e.g., https://service.run.app) from a full URL
// to use as the audience for the identity token.
func extractBaseURL(fullURL string) string {
	if strings.HasPrefix(fullURL, "http") {
		parts := strings.SplitN(fullURL, "/", 4)
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2]
		}
	}
	return fullURL
}
