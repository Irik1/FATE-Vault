package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/charactersList", charactersList)

	router.Run("localhost:8080")
}
