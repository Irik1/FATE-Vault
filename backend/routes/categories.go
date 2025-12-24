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

func CreateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var category models.CharacterCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("categories")
	_, err := coll.InsertOne(ctx, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	var category models.CharacterCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("categories")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": category}
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
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

	coll := db.Client.Database("main").Collection("categories")
	filter := bson.M{"_id": id}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

func ListCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("categories")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		c.String(http.StatusInternalServerError, "find error: %v", err)
		return
	}
	defer cur.Close(ctx)

	var results []models.CharacterCategory
	for cur.Next(ctx) {
		var category models.CharacterCategory
		if err := cur.Decode(&category); err != nil {
			c.String(http.StatusInternalServerError, "decode error: %v", err)
			return
		}

		results = append(results, category)
	}
	if err := cur.Err(); err != nil {
		c.String(http.StatusInternalServerError, "cursor error: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

