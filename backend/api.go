package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Character struct {
	ID          string   `json:"id"`
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
	byteValue, err := os.ReadFile("example.json")
	if err != nil {
		c.String(http.StatusInternalServerError, "read error: %v", err)
		return
	}

	trim := bytes.TrimSpace(byteValue)
	if len(trim) == 0 {
		c.String(http.StatusBadRequest, "empty json")
		return
	}

	var character Character
	if err := json.Unmarshal(byteValue, &character); err != nil {
		c.String(http.StatusInternalServerError, "unmarshal error: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, character)
}
