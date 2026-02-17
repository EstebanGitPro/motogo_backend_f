package middleware

import (
	"sync"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter holds a rate limiter and the last time it was accessed.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is a middleware that limits requests per IP address
// using a token bucket algorithm. Each IP gets its own bucket.
type RateLimiter struct {
	visitors sync.Map // map[string]*ipLimiter
	rps      rate.Limit
	burst    int
	done     chan struct{}
}

// NewRateLimiter creates a new rate limiter middleware.
//   - rps: requests per second allowed (e.g. 0.2 = 1 request every 5 seconds)
//   - burst: maximum burst size (allows short bursts of legitimate traffic)
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		rps:   rate.Limit(rps),
		burst: burst,
		done:  make(chan struct{}),
	}

	// Background cleanup of stale entries every 3 minutes
	go rl.cleanup(3 * time.Minute)

	return rl
}

// getLimiter retrieves or creates a rate.Limiter for the given IP.
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	if v, ok := rl.visitors.Load(ip); ok {
		entry, valid := v.(*ipLimiter)
		if !valid {
			log.Error(logger.LogMiddlewareTypeCastError, "invalid type in visitors map")
			return rate.NewLimiter(rl.rps, rl.burst)
		}
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.rps, rl.burst)
	rl.visitors.Store(ip, &ipLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	})
	return limiter
}

// cleanup removes IP entries that haven't been seen in the given expiration window.
func (rl *RateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.visitors.Range(func(key, value interface{}) bool {
				entry, ok := value.(*ipLimiter)
				if !ok {
					rl.visitors.Delete(key)
					return true
				}
				if time.Since(entry.lastSeen) > interval*2 {
					rl.visitors.Delete(key)
				}
				return true
			})
		case <-rl.done:
			return
		}
	}
}

// Stop signals the cleanup goroutine to exit.
func (rl *RateLimiter) Stop() {
	close(rl.done)
}

// Limit returns a Gin middleware handler that enforces the rate limit.
// If the client IP exceeds the allowed rate, it aborts with ErrRateLimitExceeded.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			traceID := GetRequestID(c)
			rlLog := log.WithTraceID(traceID)
			rlLog.Warn(logger.LogMiddlewareRateLimitHit,
				"client_ip", ip,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)

			_ = c.Error(domain.ErrRateLimitExceeded)
			c.Abort()
			return
		}

		c.Next()
	}
}
