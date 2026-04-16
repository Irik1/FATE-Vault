package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func webOrigin() string {
	if o := os.Getenv("WEB_ORIGIN"); o != "" {
		return o
	}
	return "http://localhost:3000"
}

// New creates and configures a new Gin engine (HTTP server).
func New() *gin.Engine {
	router := gin.Default()
	allowed := webOrigin()

	// CORS: credentials require a specific origin (not *).
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" && origin == allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowed)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, accept, origin, Cache-Control, X-Requested-With")
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
