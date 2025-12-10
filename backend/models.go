package main

import "time"

type Edition string

const (
	Core        Edition = "core"
	Accelerated Edition = "accelerated"
	Condensed   Edition = "condensed"
	Custom      Edition = "custom"
)

type Aspect struct {
	Type  string `json:"type" bson:"type"`   // e.g. "High Concept", "Trouble", or custom
	Value string `json:"value" bson:"value"` // the aspect text
}

type Refresh struct {
	Current int `json:"current" bson:"current"`
	Max     int `json:"max" bson:"max"`
}

type Consequence struct {
	Type        string `json:"type" bson:"type"`
	Size        int    `json:"size" bson:"size"`
	Description string `json:"description" bson:"description"`
	Status      string `json:"status" bson:"status"`
}

type StressBox struct {
	Size     int  `json:"size" bson:"size"`
	IsFilled bool `json:"isFilled" bson:"isFilled"`
}

type Stress struct {
	Type  string      `json:"type" bson:"type"`   // e.g. "physical", "mental"
	Boxes []StressBox `json:"boxes" bson:"boxes"` // boxes array
}

// New: SkillGroup and Stunt types (slices only)
type SkillGroup struct {
	Level  string   `json:"level" bson:"level"`
	Skills []string `json:"skills" bson:"skills"`
}

type Stunt struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type Character struct {
	ID          string   `json:"_id" bson:"_id"`
	Edition     Edition  `json:"edition" bson:"edition"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Images      []string `json:"images" bson:"images"`
	Notes       string   `json:"notes" bson:"notes"`

	Aspects      []Aspect      `json:"aspects" bson:"aspects"`
	Skills       []SkillGroup  `json:"skills" bson:"skills"`
	Refresh      Refresh       `json:"refresh" bson:"refresh"`
	Extras       string        `json:"extras" bson:"extras"`
	Stunts       []Stunt       `json:"stunts" bson:"stunts"`
	Stress       []Stress      `json:"stress" bson:"stress"`
	Consequences []Consequence `json:"consequences" bson:"consequences"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
