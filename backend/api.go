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
	ID          string   `json:"_id" bson:"_id"`
	Edition     Edition  `json:"edition" bson:"edition"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Images      []string `json:"images" bson:"images"`
	Notes       string   `json:"notes" bson:"notes"`

	Aspects      interface{}       `json:"aspects" bson:"aspects"`
	Skills       interface{}       `json:"skills" bson:"skills"`
	Refresh      int               `json:"refresh" bson:"refresh"`
	Extras       string            `json:"extras" bson:"extras"`
	Stunts       map[string]string `json:"stunts" bson:"stunts"`
	Stress       interface{}       `json:"stress" bson:"stress"`
	Consequences interface{}       `json:"consequences" bson:"consequences"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
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
