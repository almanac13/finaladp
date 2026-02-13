package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchStatus string

const (
	Scheduled MatchStatus = "scheduled"
	Live      MatchStatus = "live"
	Finished  MatchStatus = "finished"
)

type MatchEvent struct {
	Type     string `bson:"type" json:"type"` // goal | card | injury | var | sub
	TeamCode string `bson:"teamCode" json:"teamCode"`
	Minute   int    `bson:"minute" json:"minute"`

	// Common optional fields:
	Player    string `bson:"player,omitempty" json:"player,omitempty"`
	Detail    string `bson:"detail,omitempty" json:"detail,omitempty"`       // e.g. "VAR check: offside"
	CardColor string `bson:"cardColor,omitempty" json:"cardColor,omitempty"` // yellow | red

	// Substitution:
	PlayerOut string `bson:"playerOut,omitempty" json:"playerOut,omitempty"`
	PlayerIn  string `bson:"playerIn,omitempty" json:"playerIn,omitempty"`
}

type Match struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchKey string             `bson:"matchKey" json:"matchKey"` // auto-generated unique

	DateTime  time.Time   `bson:"dateTime" json:"dateTime"`
	HomeCode  string      `bson:"homeCode" json:"homeCode"`
	AwayCode  string      `bson:"awayCode" json:"awayCode"`
	HomeGoals int         `bson:"homeGoals" json:"homeGoals"`
	AwayGoals int         `bson:"awayGoals" json:"awayGoals"`
	Status    MatchStatus `bson:"status" json:"status"`

	Events []MatchEvent `bson:"events" json:"events"`
}
