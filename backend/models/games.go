package models

import "time"

type GameStatus string

const (
	InactiveGame GameStatus = "inactive"
	ActiveGame   GameStatus = "active"
	FinishedGame GameStatus = "finished"
	CanceledGame GameStatus = "canceled"
)

type GameAspect struct {
	Value         string `json:"value" bson:"value"`
	PlayerInvokes int32  `json:"playerInvokes" bson:"playerInvokes"`
	MasterInvokes int32  `json:"masterInvokes" bson:"masterInvokes"`
}

type Game struct {
	ID                  string       `json:"_id" bson:"_id"`
	Edition             Edition      `json:"edition" bson:"edition"`
	Name                string       `json:"name" bson:"name"`
	Description         string       `json:"description" bson:"description"`
	GameAspects         []GameAspect `json:"gameAspects" bson:"gameAspects"`
	LinkedCharactersIds []string     `json:"linkedCharactersIds" bson:"linkedCharactersIds"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
