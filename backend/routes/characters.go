package routes

import (
	"context"
	"net/http"
	"time"

	"FATE-Vault/backend/db"
	"FATE-Vault/backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CharactersList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("characters")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		c.String(http.StatusInternalServerError, "find error: %v", err)
		return
	}
	defer cur.Close(ctx)

	var results []map[string]interface{}
	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			c.String(http.StatusInternalServerError, "decode error: %v", err)
			return
		}

		results = append(results, doc)
	}
	if err := cur.Err(); err != nil {
		c.String(http.StatusInternalServerError, "cursor error: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func CreateCharacter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("characters")
	_, err := coll.InsertOne(ctx, character)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create character: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, character)
}

func UpdateCharacter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("characters")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": character}
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update character: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}

	c.JSON(http.StatusOK, character)
}

func DeleteCharacter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("characters")
	filter := bson.M{"_id": id}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete character: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "character deleted successfully"})
}

func FindCharacters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	filter := bson.M{}
	if edition := c.Query("edition"); edition != "" {
		filter["edition"] = edition
	}
	if name := c.Query("name"); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if characterIds := c.QueryArray("characterIds"); len(characterIds) > 0 {
		filter["_id"] = bson.M{"$in": characterIds}
	}

	coll := db.Client.Database("main").Collection("characters")
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "find error: %v", err)
		return
	}
	defer cur.Close(ctx)

	var results []models.Character
	for cur.Next(ctx) {
		var character models.Character
		if err := cur.Decode(&character); err != nil {
			c.String(http.StatusInternalServerError, "decode error: %v", err)
			return
		}

		results = append(results, character)
	}
	if err := cur.Err(); err != nil {
		c.String(http.StatusInternalServerError, "cursor error: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}
