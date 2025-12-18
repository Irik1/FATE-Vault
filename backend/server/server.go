package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

// New creates and configures a new Gin engine (HTTP server).
func New() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	registerRoutes(router)

	return router
}

// Run starts the HTTP server on the given address.
func Run(addr string) {
	if err := New().Run(addr); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
