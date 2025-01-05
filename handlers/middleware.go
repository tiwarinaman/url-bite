package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type clientData struct {
	timestamps []time.Time // Stores timestamps of requests
}

var rateLimiter = struct {
	clients map[string]*clientData
	mu      sync.Mutex
}{
	clients: make(map[string]*clientData),
}

const (
	limit       = 3                // Maximum requests per IP
	windowSize  = time.Minute      // Time window for rate limiting
	cleanupFreq = 10 * time.Minute // Frequency of cleaning old data
)

func IPRateLimiterMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	go cleanupOldEntries() // Periodic cleanup of stale entries

	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path

		// Log the incoming request's IP and path
		logger.WithFields(logrus.Fields{
			"ip":   ip,
			"path": path,
		}).Info("Request received")

		rateLimiter.mu.Lock()
		defer rateLimiter.mu.Unlock()

		now := time.Now()
		client, exists := rateLimiter.clients[ip]

		// If client does not exist, initialize data
		if !exists {
			client = &clientData{timestamps: []time.Time{now}}
			rateLimiter.clients[ip] = client
			c.Next()
			return
		}

		// Filter timestamps within the time window
		var validRequests []time.Time
		for _, ts := range client.timestamps {
			if now.Sub(ts) <= windowSize {
				validRequests = append(validRequests, ts)
			}
		}
		client.timestamps = validRequests

		// Check if the request limit has been exceeded
		if len(client.timestamps) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			// Log the rate limit exceeded event
			logger.WithFields(logrus.Fields{
				"ip":   ip,
				"path": path,
			}).Warn("Rate limit exceeded")
			return
		}

		// Add the current timestamp to the list
		client.timestamps = append(client.timestamps, now)
		c.Next()
	}
}

// Periodically clean up stale data to prevent memory leaks
func cleanupOldEntries() {
	for {
		time.Sleep(cleanupFreq)

		rateLimiter.mu.Lock()
		now := time.Now()
		for ip, client := range rateLimiter.clients {
			var validRequests []time.Time
			for _, ts := range client.timestamps {
				if now.Sub(ts) <= windowSize {
					validRequests = append(validRequests, ts)
				}
			}

			// If no valid requests left, delete the client entry
			if len(validRequests) == 0 {
				delete(rateLimiter.clients, ip)
			} else {
				client.timestamps = validRequests
			}
		}
		rateLimiter.mu.Unlock()
	}
}
