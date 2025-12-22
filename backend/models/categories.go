package models

import "time"

type CharacterCategory struct {
	ID   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`

	Subcategories []CharacterCategory `json:"subcategories,omitempty" bson:"subcategories,omitempty"`
	CharacterIDs  []string            `json:"characterIds,omitempty" bson:"characterIds,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
