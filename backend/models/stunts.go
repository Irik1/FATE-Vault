package models

import "time"

type Stunt struct {
	ID          string  `json:"_id" bson:"_id"`
	Edition     Edition `json:"edition" bson:"edition"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
