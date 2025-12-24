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

func CreateStunt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stunt models.Stunt
	if err := c.ShouldBindJSON(&stunt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("stunts")
	_, err := coll.InsertOne(ctx, stunt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stunt: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stunt)
}

func UpdateStunt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	var stunt models.Stunt
	if err := c.ShouldBindJSON(&stunt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("stunts")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": stunt}
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stunt: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "stunt not found"})
		return
	}

	c.JSON(http.StatusOK, stunt)
}

func DeleteStunt(c *gin.Context) {
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

	coll := db.Client.Database("main").Collection("stunts")
	filter := bson.M{"_id": id}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete stunt: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "stunt not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stunt deleted successfully"})
}

func ListStunts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("stunts")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		c.String(http.StatusInternalServerError, "find error: %v", err)
		return
	}
	defer cur.Close(ctx)

	var results []models.Stunt
	for cur.Next(ctx) {
		var stunt models.Stunt
		if err := cur.Decode(&stunt); err != nil {
			c.String(http.StatusInternalServerError, "decode error: %v", err)
			return
		}

		results = append(results, stunt)
	}
	if err := cur.Err(); err != nil {
		c.String(http.StatusInternalServerError, "cursor error: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

