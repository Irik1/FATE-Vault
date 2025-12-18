package server

import (
	"FATE-Vault/backend/routes"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	router.GET("/characters/list", routes.CharactersList)
	router.POST("/characters/create", routes.CreateCharacter)
	router.POST("/characters/update/:id", routes.UpdateCharacter)
	router.DELETE("/characters/delete/:id", routes.DeleteCharacter)

	router.GET("/templates", routes.GetTemplates)
}
