package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Edition string

const (
	Core        Edition = "core"
	Accelerated Edition = "accelerated"
	Condensed   Edition = "condensed"
	Custom      Edition = "custom"
)

type Character struct {
	ID          string   `json:"_id"`
	Edition     Edition  `json:"edition"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`

	Aspects      interface{} `json:"aspects"`
	Skills       interface{} `json:"skills"`
	Refresh      int         `json:"refresh"`
	Extras       string      `json:"extras"`
	Stunts       []string    `json:"stunts"`
	Stress       interface{} `json:"stress"`
	Consequences interface{} `json:"consequences"`
}

func charactersList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := mongoClient.Database("main").Collection("characters")
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
