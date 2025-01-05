package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"url-bite/utils"
)

func ShortenURL(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			OriginalURL string `json:"original_url" binding:"required,url"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			utils.LogError(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
			return
		}

		// Generate short ID
		shortID := uuid.New().String()[:8]
		createdAt := time.Now().Format(time.RFC3339)

		// Save URL mapping in the database
		_, err := db.Exec(`INSERT INTO urls (short_id, original_url, created_at) VALUES (?, ?, ?)`,
			shortID, request.OriginalURL, createdAt)
		if err != nil {
			utils.LogError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
			return
		}

		// Derive the Base URL dynamically from the request's host
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request.Host

		// Respond with the shortened URL
		c.JSON(http.StatusOK, gin.H{"short_url": baseURL + "/" + shortID})
	}
}

func RedirectURL(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortID := c.Param("shortID")

		row := db.QueryRow(`SELECT original_url FROM urls WHERE short_id = ?`, shortID)
		var originalURL string
		if err := row.Scan(&originalURL); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			} else {
				utils.LogError(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		c.Redirect(http.StatusMovedPermanently, originalURL)
	}
}
