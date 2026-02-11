package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchStatus string

const (
	Scheduled MatchStatus = "scheduled"
	Finished  MatchStatus = "finished"
)

type Goal struct {
	TeamCode string `bson:"teamCode" json:"teamCode"`
	Player   string `bson:"player" json:"player"`
	Minute   int    `bson:"minute" json:"minute"`
}

type Card struct {
	TeamCode string `bson:"teamCode" json:"teamCode"`
	Player   string `bson:"player" json:"player"`
	Color    string `bson:"color" json:"color"` // yellow | red
	Minute   int    `bson:"minute" json:"minute"`
}

type Match struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchKey string             `bson:"matchKey" json:"matchKey"` // your custom string id

	DateTime  time.Time   `bson:"dateTime" json:"dateTime"`
	HomeCode  string      `bson:"homeCode" json:"homeCode"`
	AwayCode  string      `bson:"awayCode" json:"awayCode"`
	HomeGoals int         `bson:"homeGoals" json:"homeGoals"`
	AwayGoals int         `bson:"awayGoals" json:"awayGoals"`
	Status    MatchStatus `bson:"status" json:"status"`

	Goals []Goal `bson:"goals" json:"goals"`
	Cards []Card `bson:"cards" json:"cards"`
}
