package server

import (
	"FATE-Vault/backend/routes"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	//characters
	router.GET("/characters", routes.CharactersList)
	router.POST("/characters/create", routes.CreateCharacter)
	router.POST("/characters/update/:id", routes.UpdateCharacter)
	router.DELETE("/characters/delete/:id", routes.DeleteCharacter)
	router.GET("/characters/find", routes.FindCharacters)
	router.GET("/templates", routes.GetTemplates)

	//categories
	router.GET("/categories", routes.ListCategories)
	router.POST("/categories/create", routes.CreateCategory)
	router.POST("/categories/update/:id", routes.UpdateCategory)
	router.DELETE("/categories/delete/:id", routes.DeleteCategory)
	//games
	router.GET("/games", routes.ListGames)
	router.POST("/games/create", routes.CreateGame)
	router.POST("/games/update/:id", routes.UpdateGame)
	router.DELETE("/games/delete/:id", routes.DeleteGame)

	//stunts
	router.GET("/stunts", routes.ListStunts)
	router.POST("/stunts/create", routes.CreateStunt)
	router.POST("/stunts/update/:id", routes.UpdateStunt)
	router.DELETE("/stunts/delete/:id", routes.DeleteStunt)

	//users
	router.POST("/users/register", routes.RegisterUser)
	router.POST("/users/auth", routes.AuthUser)
	router.POST("/users/update/:id", routes.AuthMiddleware(), routes.UpdateUser)
}
