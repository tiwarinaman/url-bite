package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"url-bite/config"
	"url-bite/database"
	"url-bite/handlers"
	"url-bite/utils"
)

func main() {
	// Initialize logger
	logger := utils.InitLogger()
	logger.Info("Starting URL Shortener Service...")

	// Load configurations
	config.LoadConfig()

	// Initialize SQLite database
	db, err := database.InitDB(config.AppConfig.DatabaseFile)
	if err != nil {
		logger.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	// Initialize Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/ping", handlers.Ping())
	r.HEAD("/ping", handlers.Ping()) // for light-weight request

	// Middleware for IP-based rate limiting
	r.Use(handlers.IPRateLimiterMiddleware(logger))

	// Define routes
	r.POST("/shorten", handlers.ShortenURL(db))
	r.GET("/:shortID", handlers.RedirectURL(db))

	// Start server
	logger.Fatal(r.Run(":" + config.AppConfig.ServerPort))
}
