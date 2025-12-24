package routes

import (
	"context"
	"net/http"
	"time"

	"FATE-Vault/backend/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTemplates(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	// Build filter for published templates visible to current user
	filter := visibilityFilter(c)

	coll := db.Client.Database("main").Collection("templates")
	cur, err := coll.Find(ctx, filter)
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
