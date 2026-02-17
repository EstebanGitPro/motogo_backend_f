package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRateLimitedRouter(rps float64, burst int) (*gin.Engine, *middleware.RateLimiter) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	limiter := middleware.NewRateLimiter(rps, burst)
	router.POST("/test", limiter.Limit(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return router, limiter
}

func TestRateLimiter_AllowsNormalTraffic(t *testing.T) {
	// Arrange: burst=5, so first 5 requests should pass
	router, limiter := setupRateLimitedRouter(0.2, 5)
	defer limiter.Stop()

	// Act & Assert: all burst requests should succeed
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}
}

func TestRateLimiter_BlocksExcessiveRequests(t *testing.T) {
	// Arrange: burst=3, so 4th request should be blocked
	router, limiter := setupRateLimitedRouter(0.1, 3)
	defer limiter.Stop()

	// Act: send burst+1 requests
	var lastStatus int
	for i := 0; i < 4; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		router.ServeHTTP(w, req)
		lastStatus = w.Code
	}

	// Assert: the 4th request should be blocked (200 means error middleware didn't
	// run, but we can check that the handler aborted by checking gin errors)
	// Since there's no error handler middleware in the test, the blocked request
	// returns 200 but we check differently - let's use an error-aware setup
	_ = lastStatus

	// Better approach: check that Abort was called by verifying the response
	// When c.Abort() is called without writing, status is 200 by default
	// But the handler should NOT have run, so no JSON body
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"

	// Create a router with error handler to capture the error
	gin.SetMode(gin.TestMode)
	errorRouter := gin.New()
	rl2 := middleware.NewRateLimiter(0.1, 3)
	defer rl2.Stop()

	var capturedError bool
	errorRouter.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			capturedError = true
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limited"})
		}
	})
	errorRouter.POST("/test", rl2.Limit(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Exhaust burst
	for i := 0; i < 3; i++ {
		w2 := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/test", nil)
		r.RemoteAddr = "10.0.0.1:12345"
		errorRouter.ServeHTTP(w2, r)
	}

	// This one should be rate limited
	errorRouter.ServeHTTP(w, req)
	assert.True(t, capturedError, "Rate limited request should produce an error")
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	// Arrange: burst=2
	gin.SetMode(gin.TestMode)
	router := gin.New()
	limiter := middleware.NewRateLimiter(0.1, 2)
	defer limiter.Stop()

	var handlerCalled int
	router.POST("/test", limiter.Limit(), func(c *gin.Context) {
		handlerCalled++
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act: exhaust burst for IP1
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		router.ServeHTTP(w, req)
	}

	handlerCalled = 0

	// Act: IP2 should still be allowed
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "10.0.0.2:12345"
		router.ServeHTTP(w, req)
	}

	// Assert: IP2's 2 requests should all succeed
	assert.Equal(t, 2, handlerCalled, "Different IPs should have independent rate limits")
}

func TestRateLimiter_RecoversAfterWindow(t *testing.T) {
	// Arrange: high rps so tokens refill quickly (10 rps = 1 token every 100ms)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	limiter := middleware.NewRateLimiter(10, 1) // 10 rps, burst 1
	defer limiter.Stop()

	router.POST("/test", limiter.Limit(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Exhaust burst
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/test", nil)
	req1.RemoteAddr = "10.0.0.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request should be blocked
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/test", nil)
	req2.RemoteAddr = "10.0.0.1:12345"
	router.ServeHTTP(w2, req2)

	// Wait for token to refill (at 10 rps, 150ms should be enough)
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/test", nil)
	req3.RemoteAddr = "10.0.0.1:12345"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code, "Should recover after waiting for token refill")
}

func TestRateLimiter_Stop(t *testing.T) {
	// Ensure Stop doesn't panic
	limiter := middleware.NewRateLimiter(1, 5)
	assert.NotPanics(t, func() {
		limiter.Stop()
	})
}
